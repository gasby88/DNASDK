package main

import (
	"DNA/common"
	"DNA/account"
	"DNA/core/transaction"
	"DNASDK"
	"fmt"
	"time"
)

func TransferTransaction(client *dnasdk.DnaClient, assetId common.Uint256, from, to *account.Account, amount common.Fixed64) error {
	programHash, err := client.GetAccountProgramHash(from)
	if err != nil {
		return fmt.Errorf("GetAccountProgramHash error:%s", err)
	}

	unspents, err := client.GetUnspendOutput(assetId, programHash)
	if err != nil {
		return fmt.Errorf("GetUnspendOutput error:%s", err)
	}
	if unspents == nil {
		return fmt.Errorf("GetUnspendOutput return nil")
	}

	programHashTo, err := client.GetAccountProgramHash(to)
	if err != nil {
		return fmt.Errorf("GetAccountProgramHash error:%s", err)
	}

	txInputs := make([]*transaction.UTXOTxInput, 0, 1)
	txOutputs := make([]*transaction.TxOutput, 0, 1)
	value := common.Fixed64(0)
	for _, unspent := range unspents {
		input := &transaction.UTXOTxInput{
			ReferTxID:          unspent.ReferTxID,
			ReferTxOutputIndex: unspent.ReferTxOutputIndex,
		}
		txInputs = append(txInputs, input)

		output := &transaction.TxOutput{
			AssetID:     unspent.AssetID,
			Value:       amount,
			ProgramHash: programHashTo,
		}
		txOutputs = append(txOutputs, output)

		dibs := unspent.Value - (amount - value)
		if dibs == 0 {
			value += unspent.Value
		} else if dibs > 0 {
			dibsOutput := &transaction.TxOutput{
				AssetID:     output.AssetID,
				Value:       dibs,
				ProgramHash: unspent.ProgramHash,
			}
			txOutputs = append(txOutputs, dibsOutput)
			value += unspent.Value - dibs
		}
		if value == amount {
			break
		}
	}

	if value < amount {
		return fmt.Errorf("Not enoungh utxo")
	}

	transferTx, err := client.NewTransferAssetTransaction(txInputs, txOutputs)
	if err != nil {
		return fmt.Errorf("NewTransferAssetTransaction error:%s", err)
	}

	_, err = client.SendTransaction(from, transferTx)
	if err != nil {
		return fmt.Errorf("SendTransaction error:%s", err)
	}

	_, err = client.WaitForGenerateBlock(time.Second * 30, 1)
	if err != nil {
		return fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return nil
}

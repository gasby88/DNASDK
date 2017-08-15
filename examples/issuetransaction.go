package main

import (
	"DNA/account"
	"DNA/common"
	"DNA/core/transaction"
	. "DNASDK"
	"fmt"
	"time"
)

func IssueTransaction(client *DnaClient, controller, issuer *account.Account, assetId common.Uint256, amount common.Fixed64) error{
	programHash, err := client.GetAccountProgramHash(issuer)
	if err != nil {
		return fmt.Errorf("GetProgramHash error:%s\n", err)
	}

	output := &transaction.TxOutput{
		Value:       amount,
		AssetID:     assetId,
		ProgramHash: programHash,
	}
	txOutputs := []*transaction.TxOutput{output}
	issueTx, err := client.NewIssueAssetTransaction(txOutputs)
	if err != nil {
		return fmt.Errorf("NewIssueAssetTransaction error:%s\n", err)
	}
	_, err = client.SendTransaction(controller, issueTx)
	if err != nil {
		return 	fmt.Errorf("SendTransaction error:%s\n", err)
	}

	_, err = client.WaitForGenerateBlock(time.Second * 30, 1)
	if err != nil {
		return fmt.Errorf("WaitForGenerateBlock error:%s\n", err)
	}

	return nil
}

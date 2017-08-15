package main

import (
	"DNA/account"
	"DNA/common"
	. "DNA/core/asset"
	. "DNASDK"
	"fmt"
	"time"
)

func RegisterTransaction(client *DnaClient, asset *Asset, amount common.Fixed64,issuer, controller *account.Account) (common.Uint256, error) {
	regTx, err := client.NewAssetRegisterTransaction(asset, amount, issuer, controller)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("NewAssetRegisterTransaction Asset:%+v Amount:%v Admin:%+v Account:%+v error:%s\n",
			asset,
			amount,
			issuer,
			controller,
			err)
	}

	txHash, err := client.SendTransaction(issuer, regTx)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("SendTransaction AssetRegisterTransaction error:%s\n", err)
	}

	_, err = client.WaitForGenerateBlock(time.Second * 30, 1)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s\n", err)
	}

	return txHash, nil
}

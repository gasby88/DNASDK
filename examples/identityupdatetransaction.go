package main

import (
	"DNA/account"
	"DNASDK"
	"fmt"
	"time"
)

func IdentityUpdateTransaction(client *dnasdk.DnaClient, account *account.Account, did, ddo []byte) error {
	tx, err := client.NewIdentityUpdateTransaction(account.PublicKey, did, ddo)
	if err != nil {
		return fmt.Errorf("NewIdentityUpdateTransaction error:%s", err)
	}

	_, err = client.SendTransaction(account, tx)
	if err != nil {
		fmt.Errorf("SendTransaction error:%s", err)
	}

	_, err = client.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		return fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return nil
}

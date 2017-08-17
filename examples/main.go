package main

import (
	"DNA/core/asset"
	. "DNASDK"
	"flag"
	"fmt"
	"strings"
)

var (
	DNAJsonRpcAddress string
)

func init() {
	flag.StringVar(&DNAJsonRpcAddress, "rpc", "http://localhost:20336", "The address of dna jsonrpc")
	flag.Parse()
}

func parseRpcAddress(rpcAddresses string) []string {
	return strings.Split(strings.Trim(rpcAddresses, ";"), ";")
}

func main() {
	client := NewDnaClient(parseRpcAddress(DNAJsonRpcAddress))

	walletClient := client.GetWalletClient("test")

	issuer1, err := walletClient.CreateAccount()
	if err != nil {
		fmt.Printf("Create accunt error:%s \n", err)
		return
	}
	issuer2, err := walletClient.CreateAccount()
	if err != nil {
		fmt.Printf("Create accunt error:%s \n", err)
		return
	}
	controller, err := walletClient.CreateAccount()
	if err != nil {
		fmt.Printf("Create accunt error:%s \n", err)
		return
	}

	assetName := "TS01"
	assetPrecise := byte(4)
	assetType := asset.Token
	recordType := asset.UTXO
	asset := client.CreateAsset(assetName, assetPrecise, assetType, recordType)

	assetRegAmount := client.MakeAssetAmount(20000)
	assetId, err := RegisterTransaction(client, asset, assetRegAmount, issuer1, controller)
	if err != nil {
		fmt.Printf("RegisterTransaction error:%s\n", err)
		return
	}
	fmt.Printf("RegisterTransaction success AssetId:%x Amount:%v Controller:%x\n", assetId, client.GetRawAssetAmount(assetRegAmount), controller.ProgramHash)

	assetIsuAmount := client.MakeAssetAmount(100)
	err = IssueTransaction(client, controller, issuer1, assetId, assetIsuAmount)
	if err != nil {
		fmt.Printf("IssueTransaction error:%s\n", err)
		return
	}
	fmt.Printf("IssuerTransction success Issuer:%x Amount:%v\n", issuer1.ProgramHash, client.GetRawAssetAmount(assetIsuAmount))

	assetTsfAmount := client.MakeAssetAmount(10)
	err = TransferTransaction(client, assetId, issuer1, issuer2, assetTsfAmount)
	if err != nil {
		fmt.Printf("TransferTransaction error:%s\n", err)
		return
	}
	fmt.Printf("TransferTransaction success From:%x To:%x Amount:%v\n", issuer1.ProgramHash, issuer2.ProgramHash, client.GetRawAssetAmount(assetTsfAmount))

	method := []byte("poc")
	id := []byte("123456")
	did := []byte(fmt.Sprintf("did:%s:%s", method, id))
	ddo := []byte("Hello world")
	err = SetIdentityUpdate(client, issuer1, did, ddo)
	if err != nil {
		fmt.Printf("SetIdentityUpdate error:%s\n", err)
		return
	}

	ddo2, err := GetIdentityUpdate(client, method, id)
	if err != nil {
		fmt.Printf("GetIdentityUpdate error:%s", err)
		return
	}
	if string(ddo) != string(ddo2) {
		fmt.Printf("DDO:%s not equals %s", ddo2, ddo)
		return
	}
	fmt.Printf("IdentityUpdateTransaction success DID:%s DDO:%s\n", did, ddo)
}

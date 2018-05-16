package main

import (
	"airdrop/account"
	"airdrop/config"
	"airdrop/transaction"

	"flag"
	"fmt"
)

func main() {
	var confPath string
	var mode string
	var amount int64

	// createAccount create accounts
	// sendToSub     send eth to subaddress
	// subToAirdrop  send 0 eth to airdrop address
	// withdrawToken withdraw the token from subaddress
	flag.StringVar(&mode, "m", "sendToSub", "conf path")
	flag.Int64Var(&amount, "amount", 0, "tx amount")
	flag.StringVar(&confPath, "c", "", "conf path")
	flag.Parse()
	config.Init(confPath)
	switch mode {
	case "createAccount":
		account.Create()
	case "sendToSub":
		transaction.SendToSub()
	case "sendSmartContract":
		transaction.SendSmartContract()
	case "subToAirDrop":
		transaction.SubToAirDrop()
	case "withdrawToken":
		transaction.WithdrawToken()
	default:
		fmt.Println("default")
	}
}

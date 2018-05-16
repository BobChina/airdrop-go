package transaction

import (
	"airdrop/config"
	"encoding/json"
	"fmt"
	"time"
)

// SendToSub 主账户转账到子账户
func SendToSub() {
	mainKeyStore := convertToKeystore(config.Conf.MainPrivateKey)

	// 获取当前nonce
	nonce, _ := getCurrentNonce(mainKeyStore)

	// 在console界面输出keystore
	indentMainKey, _ := json.MarshalIndent(mainKeyStore, "", "  ")
	fmt.Println(string(indentMainKey))

	addrKey := getSubAddrKey()
	for _, privateKey := range addrKey {
		subAddrkeystore := convertToKeystore(privateKey)
		receiptAddr := subAddrkeystore.Address
		fmt.Println(privateKey)        // 接收方私钥
		fmt.Println(receiptAddr.Hex()) // 接收方地址
		tx, err := SendTransaction(mainKeyStore, receiptAddr, MilliEther(config.Conf.DefaultAmount), "", nonce)
		fmt.Println(err)
		fmt.Println(tx.Hash().Hex())
		time.Sleep(1 * time.Second)
		nonce++
	}
}

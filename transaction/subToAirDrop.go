package transaction

import (
	"airdrop/config"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// SubToAirDrop 子账户send零价值eth到空投交易地址
func SubToAirDrop() {
	// 获取子钱包地址
	subAddrs := getSubAddrKey()
	for _, privateKey := range subAddrs {
		subAddrKeystore := convertToKeystore(privateKey)
		nonce, _ := getCurrentNonce(subAddrKeystore)
		tx, err := SendTransaction(subAddrKeystore, common.HexToAddress(config.Conf.AirDropAddr), common.Big0, "", nonce)
		fmt.Println(err)
		fmt.Println(tx.Hash().Hex())
	}
}

package transaction

import (
	"airdrop/config"
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const contractABI = `[{ "constant": false, "inputs": [ { "name": "_to", "type": "address" }, { "name": "_amount", "type": "uint256" } ], "name": "transfer", "outputs": [ { "name": "success", "type": "bool" } ], "payable": false, "stateMutability": "nonpayable", "type": "function" }]`

// inintContract 初始化合约
func initContract() (*bind.BoundContract, error) {
	client, err := InitClient()
	if err != nil {
		return nil, fmt.Errorf("init-client: %s", err.Error())
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, fmt.Errorf("parse-abi: %s", err.Error())
	}

	return bind.NewBoundContract(common.HexToAddress(config.Conf.AirDropAddr), parsedABI, client, client, nil), nil
}

// WithdrawToken 提取token
func WithdrawToken() (*types.Transaction, error) {
	contract, err := initContract()
	if err != nil {
		return nil, fmt.Errorf("init-contract: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Timeout)
	defer cancel()

	mainPK := convertToKeystore(config.Conf.MainPrivateKey)
	subPK := convertToKeystore(config.Conf.SubPrivateKey)

	currentNonce, _ := getCurrentNonce(subPK)

	auth := bind.NewKeyedTransactor(subPK.PrivateKey)
	auth.Nonce = big.NewInt(int64(currentNonce))
	auth.Context = ctx

	tx, err := contract.Transact(auth, "transfer", mainPK.Address, big.NewInt(1000e8))
	fmt.Println(tx.Hash().Hex())
	fmt.Println(err)
	// 查询余额
	// var output *big.Int
	// var callOpts = &bind.CallOpts{
	// 	Pending: true,
	// 	From:    mainPK.Address,
	// 	Context: ctx,
	// }
	// err = contract.Call(callOpts, &output, "balanceOf", mainPK.Address)
	// fmt.Println(err)
	// fmt.Println(output)

	return nil, nil
}

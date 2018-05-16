package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"withdraw","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"value","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_value","type":"uint256"}],"name":"burn","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"address[]"}],"name":"disableWhitelist","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"address[]"}],"name":"airdrop","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"minReq","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_value","type":"uint256"},{"name":"_minReq","type":"uint256"}],"name":"setParameters","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"finishDistribution","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"address[]"}],"name":"enableWhitelist","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"address[]"},{"name":"amounts","type":"uint256[]"}],"name":"distributeAmounts","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"getTokens","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[],"name":"distributionFinished","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"tokenAddress","type":"address"},{"name":"who","type":"address"}],"name":"getTokenBalance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"totalRemaining","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_tokenContract","type":"address"}],"name":"withdrawForeignTokens","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalDistributed","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"addresses","type":"address[]"},{"name":"amount","type":"uint256"}],"name":"distribution","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"blacklist","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_value","type":"uint256"},{"name":"_minReq","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_owner","type":"address"},{"indexed":true,"name":"_spender","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"Distr","type":"event"},{"anonymous":false,"inputs":[],"name":"DistrFinished","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"burner","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]`
)

// confJSON 配置文件格式
type confJSON struct {
	AccountURL      string        `json:"account_url"`
	URL             string        `json:"url"`
	AddrPath        string        `json:"addr_path"`
	Timeout         time.Duration `json:"timeout"`
	MainPrivateKey  string        `json:"main_private_key"`
	SubPrivateKey   string        `json:"sub_private_key"`
	AirDropAddr     string        `json:"air_drop_addr"`
	DefaultAmount   int64         `json:"default_amount"`
	DefaultGasPrice int64         `json:"default_gas_price"`
	DefaultGasLimit uint64        `json:"default_gas_limit"`
}

// initKeystore 初始化公私钥对及相应地址
func convertToKeystore(privateKeyStr string) *keystore.Key {
	// 根据字符串生成标准私钥格式
	privateKey, _ := crypto.ToECDSA((common.FromHex(privateKeyStr)))
	// 根据私钥推到出公钥
	publicKey := privateKey.PublicKey
	// 根据公钥生成地址
	address := crypto.PubkeyToAddress(publicKey)

	return &keystore.Key{
		Address:    address,
		PrivateKey: privateKey,
	}
}

var conf confJSON
var wg sync.WaitGroup

func confInit(filePath string) {
	jsonStr, _ := ioutil.ReadFile(filePath)
	json.Unmarshal(jsonStr, &conf)
}

func main() {
	var confPath string
	var mode string
	var amount int64

	// createAccount 创建账户
	// sendToSub     主账户转钱到子账户
	// subToAirdrop  子账户空投
	// subToMain     子账户将token转回
	// flag.StringVar(&mode, "m", "subToAirDrop", "conf path")
	flag.StringVar(&mode, "m", "withdrawToken", "conf path")
	flag.Int64Var(&amount, "amount", 0, "tx amount")
	flag.StringVar(&confPath, "c", "/Users/zhangan/devspace/gopath/src/airdrop/conf.json", "conf path")
	flag.Parse()
	confInit(confPath)
	switch mode {
	case "createAccount":
		for i := 0; i < 200; i++ {
			wg.Add(1)
			go CreateAccount()
		}
		wg.Wait()
	case "sendToSub":
		SendToSub()
	case "sendSmartContract":
		SendSmartContract()
	case "subToAirDrop":
		SubToAirDrop()
	case "withdrawToken":
		withdrawToken()
	default:
		fmt.Println("default")
	}
}

// SendSmartContract 发布智能合约
func SendSmartContract() {
	// mainKeyStore := convertToKeystore(mainPrivateKey)
	// SendTransaction(mainKeyStore, nil, , data string, nonce uint64)
}

// SubToAirDrop 子账户send零价值eth到空投交易地址
func SubToAirDrop() {
	// 获取子钱包地址
	subAddrs := getSubAddrKey()
	for _, privateKey := range subAddrs {
		subAddrKeystore := convertToKeystore(privateKey)
		nonce, _ := getCurrentNonce(subAddrKeystore)
		tx, err := SendTransaction(subAddrKeystore, common.HexToAddress(conf.AirDropAddr), common.Big0, "", nonce)
		fmt.Println(err)
		fmt.Println(tx.Hash().Hex())
	}
}

// CreateAccount 用于创建账户
func CreateAccount() {
	key := keystore.NewKeyForDirectICAP(rand.Reader)
	// fmt.Println(key.Id)
	fmt.Println(key.Address.Hex())
	fmt.Println(hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))
	wg.Done()
}

// SendToSub 主账户转账到子账户
func SendToSub() {
	mainKeyStore := convertToKeystore(conf.MainPrivateKey)

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
		tx, err := SendTransaction(mainKeyStore, receiptAddr, MilliEther(conf.DefaultAmount), "", nonce)
		fmt.Println(err)
		fmt.Println(tx.Hash().Hex())
		time.Sleep(1 * time.Second)
		nonce++
	}
}

func getSubAddrKey() map[string]string {
	f, err := os.Open(conf.AddrPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewScanner(f)
	i := 0
	tmp := ""
	addrKey := map[string]string{}
	for rd.Scan() {
		if i%2 == 0 {
			tmp = rd.Text()
		} else {
			addrKey[tmp] = rd.Text()
		}
		i++
	}
	return addrKey
}

var client *Client

// Client 结构体
type Client struct {
	ethclient.Client
	NetworkID *big.Int
	GasPrice  *big.Int
	GasLimit  uint64
}

// InitClient 启动ethcli
func InitClient() (*Client, error) {
	if client != nil {
		return client, nil
	}

	cli, err := ethclient.Dial(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("ethclient-dial: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()
	fmt.Println(ctx)

	networkID, err := cli.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get-chainid: %s", err.Error())
	}

	client = &Client{
		Client:    *cli,
		NetworkID: networkID,
		GasPrice:  GWei(conf.DefaultGasPrice),
		GasLimit:  conf.DefaultGasLimit,
	}
	return client, nil
}

// getCurrentNonce 获取当前交易最大nonce值
// 防止因nonce，pendding等因素导致的panic
func getCurrentNonce(from *keystore.Key) (uint64, error) {
	client, err := InitClient()
	if err != nil {
		return 0, fmt.Errorf("init-client: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	nonce, err := client.NonceAt(ctx, from.Address, nil)
	return nonce, nil
}

// SendTransaction 用于转账
// from   发送方地址
// to     接收方地址
// amount 转账额度
// data   附言
// nonce  转账次数
func SendTransaction(from *keystore.Key, to common.Address, amount *big.Int, data string, nonce uint64) (*types.Transaction, error) {
	client, err := InitClient()
	if err != nil {
		return nil, fmt.Errorf("init-client: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	tx := types.NewTransaction(nonce, to, amount, client.GasLimit, client.GasPrice, []byte(data))
	tx, err = types.SignTx(tx, types.FrontierSigner{}, from.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign-transaction: %s", err.Error())
	}

	ctx, cancel = context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("send-transaction: %s", err.Error())
	}

	fmt.Println(tx)
	return tx, nil
}

// Wei return (num * 10**exp)
func wei(num int64, exp int64) *big.Int {
	exp10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(exp), nil)
	return new(big.Int).Mul(big.NewInt(num), exp10)
}

// Wei Ethereum Base Unit
func Wei(amount int64) *big.Int {
	return wei(amount, 0)
}

// KWei Babbage = 10^3 Wei
func KWei(amount int64) *big.Int {
	return wei(amount, 3)
}

// MWei Lovelace = 10^6 Wei
func MWei(amount int64) *big.Int {
	return wei(amount, 6)
}

// GWei Shannon = 10^9 Wei
func GWei(amount int64) *big.Int {
	return wei(amount, 9)
}

// MicroEther Szabo = 10^12 Wei
func MicroEther(amount int64) *big.Int {
	return wei(amount, 12)
}

// MilliEther Finney = 10^15 Wei
func MilliEther(amount int64) *big.Int {
	return wei(amount, 15)
}

// Ether = 10^18 Wei
func Ether(amount int64) *big.Int {
	return wei(amount, 18)
}

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

	return bind.NewBoundContract(common.HexToAddress(conf.AirDropAddr), parsedABI, client, client, nil), nil
}

// withdrawToken 提取token
func withdrawToken() (*types.Transaction, error) {
	contract, err := initContract()
	if err != nil {
		return nil, fmt.Errorf("init-contract: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	mainPK := convertToKeystore(conf.MainPrivateKey)
	subPK := convertToKeystore(conf.SubPrivateKey)

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

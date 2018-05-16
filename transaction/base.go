package transaction

import (
	"airdrop/config"
	"bufio"
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

	cli, err := ethclient.Dial(config.Conf.URL)
	if err != nil {
		return nil, fmt.Errorf("ethclient-dial: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Timeout)
	defer cancel()
	fmt.Println(ctx)

	networkID, err := cli.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get-chainid: %s", err.Error())
	}

	client = &Client{
		Client:    *cli,
		NetworkID: networkID,
		GasPrice:  GWei(config.Conf.DefaultGasPrice),
		GasLimit:  config.Conf.DefaultGasLimit,
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

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Timeout)
	defer cancel()

	nonce, err := client.NonceAt(ctx, from.Address, nil)
	return nonce, nil
}

func getSubAddrKey() map[string]string {
	f, err := os.Open(config.Conf.AddrPath)
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

	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Timeout)
	defer cancel()

	tx := types.NewTransaction(nonce, to, amount, client.GasLimit, client.GasPrice, []byte(data))
	tx, err = types.SignTx(tx, types.FrontierSigner{}, from.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign-transaction: %s", err.Error())
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.Conf.Timeout)
	defer cancel()

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("send-transaction: %s", err.Error())
	}

	fmt.Println(tx)
	return tx, nil
}

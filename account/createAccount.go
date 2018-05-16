package account

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

var wg sync.WaitGroup

// Create 用于创建账户
func Create() {
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			key := keystore.NewKeyForDirectICAP(rand.Reader)
			fmt.Println(key.Address.Hex())
			fmt.Println(hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)))
			wg.Done()
		}()
	}
	wg.Wait()
}

// https://goethereumbook.org/en/client-simulated/
// https://github.com/ethereum/go-ethereum/blob/gh-pages/docs

package main

//go:generate .bin/solc -o contract --abi --bin --optimize --overwrite contract/Storage.sol
//go:generate abigen --abi contract/Storage.abi --bin contract/Storage.bin --pkg main --type Storage --out Storage.go

import (
	// "context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const key = `<<json object from keystore>>`

func main() {
	// Generate a new random account and a funded simulator
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}

	fmt.Printf("address: %s\n", address.Hex())

	blockGasLimit := uint64(4712388)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// instantiate contract
	store, err := NewStorage(common.HexToAddress("0x21e6fc92f93c8a1bb41e2be64b4e1f88a54d3576"), client)
	if err != nil {
		log.Fatalf("Failed to instantiate a Storage contract: %v", err)
	}
	// Create an authorized transactor and call the store function
	// auth, err := bind.NewStorageTransactor(strings.NewReader(key), "strong_password")
	auth2, err := bind.NewTransactor(strings.NewReader(key), "strong_password")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	// Call the store() function
	tx, err := store.Store(auth2, big.NewInt(420))
	if err != nil {
		log.Fatalf("Failed to update value: %v", err)
	}
	fmt.Printf("Update pending: 0x%x\n", tx.Hash())

	// fromAddress := auth.From
	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// value := big.NewInt(1000000000000000000) // in wei (1 eth)
	// gasLimit := uint64(21000)                // in units
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	// var data []byte
	// tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	// signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex()) // tx sent: 0xec3ceb05642c61d33fa6c951b54080d1953ac8227be81e7b5e4e2cfed69eeb51

	// client.Commit()

	// receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if receipt == nil {
	// 	log.Fatal("receipt is nil. Forgot to commit?")
	// }

	// fmt.Printf("status: %v\n", receipt.Status) // status: 1
}

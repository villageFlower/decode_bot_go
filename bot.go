package main

import (
	"context"
	"fmt"
	"log"
	"simon/decode"
	"strings"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"io/ioutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"os"
)

func getCurrentClient() *gethclient.Client {
	rpcClient, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/9110a9490de6477184406113ce4854a4")
	if err != nil {
		log.Fatal(err)
	}
	client := gethclient.New(rpcClient)
	return client
}

func main() {
	const rawUrl = "wss://mainnet.infura.io/ws/v3/9110a9490de6477184406113ce4854a4"
	const txChCap = 30
	const contractAddress = "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B"

	client := getCurrentClient()

	txCh := make(chan *types.Transaction, txChCap)
	sub, err := client.SubscribeFullPendingTransactions(context.Background(), txCh)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// load contract ABI

	contactAbi, err := os.Open("./decode/abi.json")
	if err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(contactAbi)
	defer contactAbi.Close()
	abi, err := abi.JSON(strings.NewReader(string(out)))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case tx := <-txCh:
			// fmt.Printf("New tx: %v\n", tx.To())
			if tx.To() != nil {
				if tx.To().String() == contractAddress {
					fmt.Printf("Swap found, decoding has: %v\n", tx.Hash())

					decode.DecodeContract(abi, tx.Data())

				}
			}

		case err := <-sub.Err():
			log.Fatal(err)
		}
	}
}

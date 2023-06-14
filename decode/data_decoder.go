package decode

import (
	"math/big"
	"fmt"
	"time"
)

type decodeData struct {
	recepient []byte
	amountIn *big.Int
	amountOutMin *big.Int
	path [][]byte
 }

func DecodeInput(tx []byte, start time.Time) {
	result := decodeData{}
	result.recepient = make([]byte, 20)
	copy(result.recepient, tx[12:32])
	result.amountIn = new(big.Int).SetBytes(tx[32:64])
	result.amountOutMin = new(big.Int).SetBytes(tx[64:96])

	pathCount := new(big.Int).Div(new(big.Int).SetBytes(tx[97:128]), new(big.Int).SetInt64(32))

	result.path = make([][]byte, pathCount.Int64())
	for i := int64(0); i < pathCount.Int64()-1; i++ {
		result.path[i] = make([]byte, 32)
		copy(result.path[i], tx[128+i*32:128+(i+1)*32])
	}



	fmt.Printf("amountIn: %d\n", result.amountIn)
	fmt.Printf("recepient: %x\n", result.recepient)
	fmt.Printf("amountOutMin: %d\n", result.amountOutMin)
	fmt.Printf("path: %x\n", result.path)

	fmt.Println(" ")

	elapsed := time.Since(start)
	fmt.Printf("decode took %s\n\n", elapsed)
	fmt.Println(" ")

}



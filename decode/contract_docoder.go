package decode

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"os"
	"log"
	"io/ioutil"
	"time"
)

// refer
// https://github.com/ethereum/web3.py/blob/master/web3/contract.py#L435
func decodeTransactionInputData(contractABI *abi.ABI, data []byte, start time.Time) map[string]interface{} {
	methodSigData := data[:4]
	inputsSigData := data[4:]
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		log.Fatal(err)
	}
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		log.Fatal(err)
	} 
	commands := inputsMap["commands"].([]byte)

	if commands[len(commands)-1] == 0x08 {
		fmt.Println("08 command detected, decoding....")
		data := inputsMap["inputs"].([][]byte)
		DecodeInput(data[len(data)-1], start)
	}else{
		fmt.Println("invalid command, skiped")
		fmt.Println(" ")
	}

	return inputsMap
}

func DecodeContract(txInput []byte) (result map[string]interface{}) {
	// load contract ABI
	start := time.Now()
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
	result = decodeTransactionInputData(&abi, txInput, start)
	return result
}
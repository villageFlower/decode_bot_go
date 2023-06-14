package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
    "os"
)

type (
	RawABIResponse struct {
		Status  *string `json:"status"`
		Message *string `json:"message"`
		Result  *string `json:"result"`
	}
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func GetContractRawABI() (*RawABIResponse, error) {
	client := resty.New()
	rawABIResponse := &RawABIResponse{}
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"module":  "contract",
			"action":  "getabi",
			"address": "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B",
			"apikey":  "35248DCBB9EQA84WG1FSSC3376X2DMX2NZ",
		}).
		SetResult(rawABIResponse).
		Get("https://api.etherscan.io/api")

	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf(fmt.Sprintf("Get contract raw abi failed: %s\n", resp))
	}
	if *rawABIResponse.Status != "1" {
		return nil, fmt.Errorf(fmt.Sprintf("Get contract raw abi failed: %s\n", *rawABIResponse.Result))
	}

	return rawABIResponse, nil
}

func main() {


    f, err := os.Create("abi.json")
    check(err)

    defer f.Close()

    d2, err := GetContractRawABI()
    n2, err := f.Write([]byte(*d2.Result))
	fmt.Printf("wrote %d bytes\n", n2)
    check(err)


    f.Sync()

}
package main

import (
	"fmt"
	"bsnBaseCC/bsnchaincode"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	err := shim.Start(new(bsnchaincode.BsnChainCode))
	if err != nil {
		fmt.Printf("Error starting BsnChainCode: %s", err)
	}
}

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"airBlock/chaincode"
)

func main() {
	docChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("ERROR: %v", err)
	}
	if err := docChaincode.Start(); err != nil {
		log.Panicf("ERROR: %v", err)
	}
}

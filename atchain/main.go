package main

import (
	"log"

	cc "atchain/contract"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	contract, err := contractapi.NewChaincode(new(cc.SmartContract))
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	err = contract.Start()
	if err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}

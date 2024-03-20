package main

import (
	"bootcamp/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Erro ao criar chaincode 'bootcamp': %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Erro ao iniciar chaincode 'bootcamp': %v", err)
	}
}

/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	apiGoContract := new(ApiGoContract)
	apiGoContract.Info.Version = "0.0.1"
	apiGoContract.Info.Description = "My Smart Contract"
	apiGoContract.Info.License = new(metadata.LicenseMetadata)
	apiGoContract.Info.License.Name = "Apache-2.0"
	apiGoContract.Info.Contact = new(metadata.ContactMetadata)
	apiGoContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(apiGoContract)
	chaincode.Info.Title = "APIGo chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from ApiGoContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}

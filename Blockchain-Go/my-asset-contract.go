/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	startKey = "000"
	endKey = "999"
)

// MyAssetContract contract for managing CRUD for MyAsset
type MyAssetContract struct {
	contractapi.Contract
}

// MyAssetExists returns true when asset with given ID exists in world state
func (c *MyAssetContract) MyAssetExists(ctx contractapi.TransactionContextInterface, myAssetID string) (bool, error) {
	data, err := ctx.GetStub().GetState(myAssetID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateMyAsset creates a new instance of MyAsset
func (c *MyAssetContract) CreateMyAsset(ctx contractapi.TransactionContextInterface, myAssetID string, value string) error {
	exists, err := c.MyAssetExists(ctx, myAssetID)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("the asset %s already exists", myAssetID)
	}

	myAsset := new(MyAsset)
	myAsset.Value = value

	bytes, _ := json.Marshal(myAsset)

	return ctx.GetStub().PutState(myAssetID, bytes)
}

// ReadMyAsset retrieves an instance of MyAsset from the world state
func (c *MyAssetContract) ReadMyAsset(ctx contractapi.TransactionContextInterface, myAssetID string) (*MyAsset, error) {
	exists, err := c.MyAssetExists(ctx, myAssetID)
	if err != nil {
		return nil, fmt.Errorf("cCould not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", myAssetID)
	}

	bytes, _ := ctx.GetStub().GetState(myAssetID)

	myAsset := new(MyAsset)

	err = json.Unmarshal(bytes, myAsset)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type MyAsset")
	}

	return myAsset, nil
}

// UpdateMyAsset retrieves an instance of MyAsset from the world state and updates its value
func (c *MyAssetContract) UpdateMyAsset(ctx contractapi.TransactionContextInterface, myAssetID string, newValue string) error {
	exists, err := c.MyAssetExists(ctx, myAssetID)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("the asset %s does not exist", myAssetID)
	}

	myAsset := new(MyAsset)
	myAsset.Value = newValue

	bytes, _ := json.Marshal(myAsset)

	return ctx.GetStub().PutState(myAssetID, bytes)
}

// DeleteMyAsset deletes an instance of MyAsset from the world state
func (c *MyAssetContract) DeleteMyAsset(ctx contractapi.TransactionContextInterface, myAssetID string) error {
	exists, err := c.MyAssetExists(ctx, myAssetID)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("the asset %s does not exist", myAssetID)
	}

	return ctx.GetStub().DelState(myAssetID)
}

// queryAllAssets retrive all instances from the world state @Transaction(false)
func (c *MyAssetContract) QueryAllAssets(ctx contractapi.TransactionContextInterface) ([]MyAsset, error) {
	var listaDeAssets []MyAsset
	it, errStub := ctx.GetStub().GetStateByRange(startKey,endKey)
	if errStub != nil {
		return nil, fmt.Errorf("could stablish the range from world state. %s", errStub)
	}
	defer it.Close()
	for it.HasNext() {
		res, errIt := it.Next()
		if errIt != nil {
			return nil, errIt
		}
		var miAsset MyAsset
		errJson := json.Unmarshal(res.Value, &miAsset)
		if errJson != nil {
			return nil, errJson
		}
		listaDeAssets = append(listaDeAssets, miAsset)
	}
	return listaDeAssets, nil
}

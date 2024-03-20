/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ApiGoContract contract for managing CRUD for ApiGo
type ApiGoContract struct {
	contractapi.Contract
}

// ApiGoExists returns true when asset with given ID exists in world state
func (c *ApiGoContract) ApiGoExists(ctx contractapi.TransactionContextInterface, apiGoID string) (bool, error) {
	data, err := ctx.GetStub().GetState(apiGoID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateApiGo creates a new instance of ApiGo
func (c *ApiGoContract) CreateApiGo(ctx contractapi.TransactionContextInterface, apiGoID string, value string) error {
	exists, err := c.ApiGoExists(ctx, apiGoID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", apiGoID)
	}

	apiGo := new(ApiGo)
	apiGo.Value = value

	bytes, _ := json.Marshal(apiGo)

	return ctx.GetStub().PutState(apiGoID, bytes)
}

// ReadApiGo retrieves an instance of ApiGo from the world state
func (c *ApiGoContract) ReadApiGo(ctx contractapi.TransactionContextInterface, apiGoID string) (*ApiGo, error) {
	exists, err := c.ApiGoExists(ctx, apiGoID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", apiGoID)
	}

	bytes, _ := ctx.GetStub().GetState(apiGoID)

	apiGo := new(ApiGo)

	err = json.Unmarshal(bytes, apiGo)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type ApiGo")
	}

	return apiGo, nil
}

// UpdateApiGo retrieves an instance of ApiGo from the world state and updates its value
func (c *ApiGoContract) UpdateApiGo(ctx contractapi.TransactionContextInterface, apiGoID string, newValue string) error {
	exists, err := c.ApiGoExists(ctx, apiGoID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", apiGoID)
	}

	apiGo := new(ApiGo)
	apiGo.Value = newValue

	bytes, _ := json.Marshal(apiGo)

	return ctx.GetStub().PutState(apiGoID, bytes)
}

// DeleteApiGo deletes an instance of ApiGo from the world state
func (c *ApiGoContract) DeleteApiGo(ctx contractapi.TransactionContextInterface, apiGoID string) error {
	exists, err := c.ApiGoExists(ctx, apiGoID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", apiGoID)
	}

	return ctx.GetStub().DelState(apiGoID)
}

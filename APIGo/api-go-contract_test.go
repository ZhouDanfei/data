/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testApiGo := new(ApiGo)
	testApiGo.Value = "set value"
	apiGoBytes, _ := json.Marshal(testApiGo)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "apiGokey").Return(apiGoBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestApiGoExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(ApiGoContract)

	exists, err = c.ApiGoExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.ApiGoExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.ApiGoExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateApiGo(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(ApiGoContract)

	err = c.CreateApiGo(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateApiGo(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateApiGo(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadApiGo(t *testing.T) {
	var apiGo *ApiGo
	var err error

	ctx, _ := configureStub()
	c := new(ApiGoContract)

	apiGo, err = c.ReadApiGo(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, apiGo, "should not return ApiGo when exists errors when reading")

	apiGo, err = c.ReadApiGo(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, apiGo, "should not return ApiGo when key does not exist in world state when reading")

	apiGo, err = c.ReadApiGo(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type ApiGo", "should error when data in key is not ApiGo")
	assert.Nil(t, apiGo, "should not return ApiGo when data in key is not of type ApiGo")

	apiGo, err = c.ReadApiGo(ctx, "apiGokey")
	expectedApiGo := new(ApiGo)
	expectedApiGo.Value = "set value"
	assert.Nil(t, err, "should not return error when ApiGo exists in world state when reading")
	assert.Equal(t, expectedApiGo, apiGo, "should return deserialized ApiGo from world state")
}

func TestUpdateApiGo(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(ApiGoContract)

	err = c.UpdateApiGo(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateApiGo(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateApiGo(ctx, "apiGokey", "new value")
	expectedApiGo := new(ApiGo)
	expectedApiGo.Value = "new value"
	expectedApiGoBytes, _ := json.Marshal(expectedApiGo)
	assert.Nil(t, err, "should not return error when ApiGo exists in world state when updating")
	stub.AssertCalled(t, "PutState", "apiGokey", expectedApiGoBytes)
}

func TestDeleteApiGo(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(ApiGoContract)

	err = c.DeleteApiGo(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteApiGo(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteApiGo(ctx, "apiGokey")
	assert.Nil(t, err, "should not return error when ApiGo exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "apiGokey")
}

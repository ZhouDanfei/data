package assets

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func GetState(ctx contractapi.TransactionContextInterface, key string) ([]byte, error) {
	js, err := ctx.GetStub().GetState(key)
	if js == nil {
		return nil, errors.New(key + " not found")
	}
	return js, err
}

func GetAll(ctx contractapi.TransactionContextInterface, startKey string) ([][]byte, error) {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	lst := make([][]byte, 0)
	defer iterator.Close()
	for iterator.HasNext() {
		i, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		if startsWith(startKey, i.Key) {
			lst = append(lst, i.Value)
		}
	}
	return lst, err
}

func startsWith(pattern string, value string) bool {
	if len(pattern) > len(value) {
		return false
	}
	return value[:len(pattern)] == pattern
}

func RandUint() int {
	return RandUintn(math.MaxInt32)
}

func RandUintn(n int) int {
	return rand.New(rand.NewSource(time.Now().Unix())).Intn(n)
}

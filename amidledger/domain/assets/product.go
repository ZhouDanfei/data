package assets

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Product struct {
	Name     string  `json:"name"`
	Producer string  `json:"producer"`
	Date     string  `json:"date"`
	Safetime int     `json:"safetime"`
	Temp1    int     `json:"temp1"`
	Temp2    int     `json:"temp2"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	Shop     string  `json:"shop"`
}

func GetProduct(ctx contractapi.TransactionContextInterface, name string) (*Product, error) {
	js, err := GetState(ctx, "product@"+name)
	if err != nil {
		return nil, err
	}
	var product Product
	if err := json.Unmarshal(js, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	productsbt, err := GetAll(ctx, "product@")
	if err != nil {
		return nil, err
	}
	products := make([]*Product, 0, len(productsbt))
	for _, js := range productsbt {
		var product Product
		if err := json.Unmarshal(js, &product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func SetProduct(ctx contractapi.TransactionContextInterface, name string, product *Product) error {
	js, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("product@"+name, js)
}

const ProviderLogin = "goldfish"

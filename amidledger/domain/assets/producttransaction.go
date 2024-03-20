package assets

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ProductReturnTx struct {
	Id          string `json:"id"`
	ProductName string `json:"productname"`
	UserLogin   string `json:"userlogin"`
	ShopNumber  string `json:"shopnumber"`
}

func SetProductTx(ctx contractapi.TransactionContextInterface, id string, productTx *ProductReturnTx) error {
	js, err := json.Marshal(productTx)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("producttx@"+id, js)
}

func GetProductTx(ctx contractapi.TransactionContextInterface, id string) (*ProductReturnTx, error) {
	js, err := GetState(ctx, "producttx@"+id)
	if err != nil {
		return nil, err
	}
	var productTx ProductReturnTx

	if err := json.Unmarshal(js, &productTx); err != nil {
		return nil, err
	}
	return &productTx, nil
}

func GetAllProductTx(ctx contractapi.TransactionContextInterface) ([]*ProductReturnTx, error) {
	productTxbt, err := GetAll(ctx, "producttx@")
	if err != nil {
		return nil, err
	}
	productTxList := make([]*ProductReturnTx, 0, len(productTxbt))
	for _, js := range productTxbt {
		var product ProductReturnTx
		if err := json.Unmarshal(js, &product); err != nil {
			return nil, err
		}
		productTxList = append(productTxList, &product)
	}
	return productTxList, nil
}

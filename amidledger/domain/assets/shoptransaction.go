package assets

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ShopTransaction struct {
	Id             string   `json:"id"`
	ShopNumber     string   `json:"shopnumber"`
	StartCost      float64  `json:"startcost"`
	CurrentCost    float64  `json:"currentcost"`
	Temp           []int    `json:"temp"`
	Product        *Product `json:"product"`
	TempFailAmount uint     `json:"tempfailamount"`
}

func GetShopTx(ctx contractapi.TransactionContextInterface, id string) (*ShopTransaction, error) {
	js, err := GetState(ctx, "shoptx@"+id)
	if err != nil {
		return nil, err
	}
	var shopTx ShopTransaction
	err = json.Unmarshal(js, &shopTx)
	return &shopTx, err
}

func SetShopTx(ctx contractapi.TransactionContextInterface, id string, shopTx *ShopTransaction) error {
	js, err := json.Marshal(shopTx)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("shoptx@"+id, js)
}

func GetAllShopTx(ctx contractapi.TransactionContextInterface) ([]*ShopTransaction, error) {
	shopTxbt, err := GetAll(ctx, "shoptx@")
	if err != nil {
		return nil, err
	}
	shopTxList := make([]*ShopTransaction, 0, len(shopTxbt))
	for _, js := range shopTxbt {
		var shopTx ShopTransaction
		fmt.Println(string(js), js)
		if err = json.Unmarshal(js, &shopTx); err != nil {
			return nil, err
		}
		shopTxList = append(shopTxList, &shopTx)
	}
	return shopTxList, err
}

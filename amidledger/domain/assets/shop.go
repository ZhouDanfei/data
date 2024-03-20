package assets

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Shop struct {
	Number   string     `json:"number"`
	City     string     `json:"city"`
	Products []*Product `json:"products"`
}

func GetShop(ctx contractapi.TransactionContextInterface, number string) (*Shop, error) {
	js, err := GetState(ctx, "shop@"+number)
	if err != nil {
		return nil, err
	}
	var shop Shop
	if err := json.Unmarshal(js, &shop); err != nil {
		return nil, err
	}
	return &shop, nil
}

func SetShop(ctx contractapi.TransactionContextInterface, number string, shop *Shop) error {
	js, err := json.Marshal(shop)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("shop@"+number, js)
}

func GetAllShops(ctx contractapi.TransactionContextInterface) ([]*Shop, error) {
	shopsbt, err := GetAll(ctx, "shop@")
	if err != nil {
		return nil, err
	}
	shops := make([]*Shop, 0, len(shopsbt))
	for _, js := range shopsbt {
		var shop Shop
		json.Unmarshal(js, &shop)
		shops = append(shops, &shop)
	}
	return shops, err
}

func InitShops(ctx contractapi.TransactionContextInterface) error {
	if err := SetShop(ctx, "1", &Shop{"1", "Дмитров", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetUser(ctx, "1", &FabricUser{"1", "1", "shop", 1000, "1", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetShop(ctx, "2", &Shop{"2", "Калуга", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetUser(ctx, "2", &FabricUser{"2", "2", "shop", 900, "2", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetShop(ctx, "3", &Shop{"3", "Москва", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetUser(ctx, "3", &FabricUser{"3", "3", "shop", 1050, "3", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetShop(ctx, "4", &Shop{"4", "Рязань", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetUser(ctx, "4", &FabricUser{"4", "4", "shop", 1000, "4", make([]*Product, 0)}); err != nil {
		return err
	}
	return nil

}

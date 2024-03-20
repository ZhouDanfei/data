package assets

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type FabricUser struct {
	Login    string     `json:"login"`
	Password string     `json:"password"`
	Role     string     `json:"role"`
	Balance  float64    `json:"balance"`
	Fio      string     `json:"fio"`
	Products []*Product `json:"products"`
}

func GetUser(ctx contractapi.TransactionContextInterface, login string) (*FabricUser, error) {
	js, err := GetState(ctx, "user@"+login)
	if err != nil {
		return nil, err
	}
	var user FabricUser
	if err := json.Unmarshal(js, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*FabricUser, error) {
	usersbt, err := GetAll(ctx, "user@")
	if err != nil {
		return nil, err
	}
	users := make([]*FabricUser, 0, len(usersbt))

	for _, js := range usersbt {
		var user FabricUser
		if err := json.Unmarshal(js, &user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func SetUser(ctx contractapi.TransactionContextInterface, login string, user *FabricUser) error {
	js, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState("user@"+login, js)
}

func InitUsers(ctx contractapi.TransactionContextInterface) error {
	if err := SetUser(ctx, ProviderLogin, &FabricUser{ProviderLogin, ProviderLogin, "provider", 100, "Золотая рыбка", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetUser(ctx, "roman", &FabricUser{"roman", "roman", "customer", 80, "Романов Роман Романович", make([]*Product, 0)}); err != nil {
		return err
	}
	if err := SetUser(ctx, "nikola", &FabricUser{"nikola", "nikola", "customer", 90, "Николаев Николай Николаевич", make([]*Product, 0)}); err != nil {
		return err
	}
	return nil
}

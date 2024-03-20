package smartcontract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// 建立Asset結構，後面為Marshal成json檔後所呈現的方式
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           string `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue string `json:"appraisedValue"`
}

// 初始化帳本資料
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset3", Color: "blue", Size: "5", Owner: "Tomoko", AppraisedValue: "300"},
		{ID: "asset4", Color: "red", Size: "5", Owner: "Brad", AppraisedValue: "400"},
		{ID: "asset5", Color: "green", Size: "10", Owner: "Jin Soo", AppraisedValue: "500"},
		{ID: "asset6", Color: "yellow", Size: "10", Owner: "Max", AppraisedValue: "600"},
		{ID: "asset7", Color: "black", Size: "15", Owner: "Adriana", AppraisedValue: "700"},
		{ID: "asset8", Color: "white", Size: "15", Owner: "Michel", AppraisedValue: "800"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// 建立資料
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size string, owner string, appraisedValue string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// 讀取所輸入ID的Asset
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// 更新已儲存的Asset
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size string, owner string, appraisedValue string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// 將原始資料覆蓋
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// 將所輸入ID的Asset從world state中刪除
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// 當輸入ID的資料存在時返回True
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// 交易Asset
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// 取得所有Asset
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "") //抓取一個範圍的資料，參數為key的範圍，回傳後為iterator的型態
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close() //close掉iterator

	var assets []*Asset //宣告空array
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next() //把下一個iterator的response抓回來
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset) //透過json.Unmarshal轉成user的structure
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset) //append到user的陣列
	} //當for迴圈跑完，所有資料會被append完

	return assets, nil
}

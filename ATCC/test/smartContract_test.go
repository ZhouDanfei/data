package test

import (
	"atcc/smartcontract"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/longbridgeapp/assert"
)

var Stub *shimtest.MockStub
var Scc *contractapi.ContractChaincode

//先定義asset資料方便測試
var asset1 smartcontract.Asset = smartcontract.Asset{
	ID:             "1",
	Color:          "blue",
	Size:           "5",
	Owner:          "Greg",
	AppraisedValue: "300",
}
var asset2 smartcontract.Asset = smartcontract.Asset{
	ID:             "2",
	Color:          "red",
	Size:           "7",
	Owner:          "SHEN",
	AppraisedValue: "500",
}

func TestMain(m *testing.M) { //建立測試執行程式
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	log.SetOutput(ioutil.Discard)
}

func NewStub() { //建立新的智能合約
	Scc, err := contractapi.NewChaincode(new(smartcontract.SmartContract))
	if err != nil {
		log.Println("NewChaincode failed", err)
		os.Exit(0)
	}
	Stub = shimtest.NewMockStub("main", Scc)
}

//Test Functions

func Test_CreateAsset(t *testing.T) {
	fmt.Println("Test_CreateAsset-----------------")
	NewStub()

	err := MockCreateAsset(asset1.ID, asset1.Color, asset1.Size, asset1.Owner, asset1.AppraisedValue) //建立Asset
	//若有錯誤則中止
	if err != nil {
		t.FailNow()
	}
}

func Test_AssetExists(t *testing.T) {
	fmt.Println("Test_AssetExists-----------------")
	NewStub()

	err := MockCreateAsset(asset1.ID, asset1.Color, asset1.Size, asset1.Owner, asset1.AppraisedValue)
	if err != nil {
		t.FailNow()
	}

	result, err := MockAssetExists(asset1.ID) //查詢Asset是否存在
	if err != nil {
		t.FailNow()
	}
	fmt.Println("result: ", result)
	assert.Equal(t, result, true)
}

func Test_ReadAsset(t *testing.T) {
	fmt.Println("Test_ReadAsset-----------------")
	NewStub()

	err := MockCreateAsset(asset1.ID, asset1.Color, asset1.Size, asset1.Owner, asset1.AppraisedValue)
	if err != nil {
		t.FailNow()
	}

	assetJson, err := MockReadAsset(asset1.ID)
	if err != nil {
		fmt.Println("read Asset error", err)
	}
	fmt.Println("assetJson: ", assetJson)
	assert.Equal(t, assetJson.ID, asset1.ID)
	assert.Equal(t, assetJson.Color, asset1.Color)
	assert.Equal(t, assetJson.Size, asset1.Size)
	assert.Equal(t, assetJson.Owner, asset1.Owner)
	assert.Equal(t, assetJson.AppraisedValue, asset1.AppraisedValue)
}

func Test_UpdateAsset(t *testing.T) {
	fmt.Println("Test_UpdateAsset-----------------")
	NewStub()

	err := MockCreateAsset(asset1.ID, asset1.Color, asset1.Size, asset1.Owner, asset1.AppraisedValue)
	if err != nil {
		t.FailNow()
	}
	//update key=asset1.ID的資料，以change color與change size等模擬更改資料
	MockUpdateAsset(asset1.ID, "change color", "change size", "change owner", "change appraisedValue")

	//取得asset1.ID資料
	assetJson, err := MockReadAsset(asset1.ID)
	//錯誤則印出
	if err != nil {
		fmt.Println("read Asset", err)
	}
	fmt.Println("assetJson: ", assetJson)
	assert.Equal(t, assetJson.ID, asset1.ID)
	assert.Equal(t, assetJson.Color, "change color")
	assert.Equal(t, assetJson.Size, "change size")
	assert.Equal(t, assetJson.Owner, "change owner")
	assert.Equal(t, assetJson.AppraisedValue, "change appraisedValue")

}

func Test_DeleteAsset(t *testing.T) {
	fmt.Println("Test_DeleteAsset-----------------")
	NewStub()
	//asset1.ID是稍後要刪除的資料
	err := MockCreateAsset(asset1.ID, asset1.Color, asset1.Size, asset1.Owner, asset1.AppraisedValue)
	if err != nil {
		//若有錯誤則中斷
		t.FailNow()
	}
	//刪除
	MockDeleteAsset(asset1.ID)
	//取asset1.ID資料
	assetsJson, err := MockReadAsset(asset1.ID)
	//ReadAsset如果
	if err != nil {
		fmt.Println("read Asset", err)
	}
	fmt.Println(assetsJson)
	assert.Equal(t, err, errors.New("ReadAsset error"))
}

func Test_GetAllAssets(t *testing.T) {
	fmt.Println("MockGetAllAssets-----------------")
	NewStub()

	MockCreateAsset(asset1.ID, asset1.Color, asset1.Size, asset1.Owner, asset1.AppraisedValue)
	MockCreateAsset(asset2.ID, asset2.Color, asset2.Size, asset2.Owner, asset2.AppraisedValue)

	assets, err := MockGetAllAssets()
	if err != nil {
		fmt.Println("GetAllAssets error", err)
	}
	fmt.Println("Assets: ", assets)
	assert.Equal(t, len(assets), 2) //檢查長度是否為2
}

//Mock Functions

func MockCreateAsset(id string, color string, size string, owner string, appraisedValue string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("CreateAsset"),
			[]byte(id),
			[]byte(color),
			[]byte(size),
			[]byte(owner),
			[]byte(appraisedValue),
		})

	if res.Status != shim.OK {
		fmt.Println("CreateAsset failed", string(res.Message))
		return errors.New("CreateAsset error")
	}
	return nil
}

func MockAssetExists(id string) (bool, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("AssetExists"), []byte(id)})
	if res.Status != shim.OK {
		return false, errors.New("AssetExists error")
	}
	var result bool = false
	json.Unmarshal(res.Payload, &result)
	return result, nil
}

func MockReadAsset(id string) (*smartcontract.Asset, error) {
	var result smartcontract.Asset
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("ReadAsset"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("ReadAsset failed", string(res.Message))
		return nil, errors.New("ReadAsset error")
	}
	json.Unmarshal(res.Payload, &result)
	return &result, nil
}

func MockUpdateAsset(id string, color string, size string, owner string, appraisedValue string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("UpdateAsset"),
			[]byte(id),
			[]byte(color),
			[]byte(size),
			[]byte(owner),
			[]byte(appraisedValue),
		})
	if res.Status != shim.OK {
		fmt.Println("UpdateAsset failed", string(res.Message))
		return errors.New("UpdateAsset error")
	}
	return nil
}

func MockDeleteAsset(id string) error {
	res := Stub.MockInvoke("uuid",
		[][]byte{
			[]byte("DeleteAsset"),
			[]byte(id),
		})
	if res.Status != shim.OK {
		fmt.Println("DeleteAsset failed", string(res.Message))
		return errors.New("DeleteAsset error")
	}
	return nil
}

func MockGetAllAssets() ([]*smartcontract.Asset, error) {
	res := Stub.MockInvoke("uuid", [][]byte{[]byte("GetAllAssets")})
	if res.Status != shim.OK {
		fmt.Println("GetAllAssets failed", string(res.Message))
		return nil, errors.New("GetAllAssets error")
	}
	var assets []*smartcontract.Asset
	json.Unmarshal(res.Payload, &assets)
	return assets, nil
}

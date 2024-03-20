package chaincode

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Doc struct {
	Type         string   `json:"type"`
	Registration string   `json:"registration"`
	Date         string   `json:"date"`
	Text         string   `json:"text"`
	Hash         [32]byte `json:"hash"`
}

func (s *SmartContract) CreateDoc(
	ctx contractapi.TransactionContextInterface,
	docType,
	registration,
	date,
	text string,
) error {
	doc := Doc{
		Type:         docType,
		Registration: registration,
		Date:         date,
		Text:         text,
		Hash:         sha256.Sum256([]byte(text)),
	}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(docType, docJSON)
}

func (s *SmartContract) GetAllDocs(ctx contractapi.TransactionContextInterface) ([]*Doc, error) {
	iter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	var docs []*Doc
	for iter.HasNext() {
		result, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var doc Doc
		err = json.Unmarshal(result.Value, &doc)
		if err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}

package utils

import (
	"bootcamp/chaincode/messages"
	"bytes"
	"encoding/json"
	"fmt"
	s "strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ParseJSON efetua parsing de string dataJSON para um type de interface v.
func ParseJSON(dataJSON string, v interface{}) error {
	dataBytes := []byte(dataJSON)
	jsonValido := json.Valid(dataBytes)
	if !jsonValido {
		return fmt.Errorf(messages.JSONInvalido)
	}

	dec := json.NewDecoder(bytes.NewReader(dataBytes))
	dec.DisallowUnknownFields()
	err := dec.Decode(&v)
	if err != nil {
		return err
	}

	return nil
}

// ValidateMSPID valida se uma determinada string msp está como prefixo no MSPID do contexto.
func ValidateMSPID(ctx contractapi.TransactionContextInterface, msp string) error {
	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	if s.HasPrefix(s.ToUpper(mspID), s.ToUpper(msp)) {
		return nil
	}

	return fmt.Errorf(messages.TransacaoRestritaOrgao, s.ToUpper(msp))
}

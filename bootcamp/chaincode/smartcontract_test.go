package chaincode_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"bootcamp/chaincode"
	"bootcamp/chaincode/messages"
	"bootcamp/chaincode/mocks"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

const pessoafisicaJSON string = `{
	"cpf": "12345",
	"nomeCompleto": "JOÃO DA SILVA",
	"dataNascimento": "1980-03-01"
}`

const cnhJSON string = `{
	"numero": "987",
	"dataValidade": "2030-01-01"
}`

var chaincodeStubInstance = &mocks.ChaincodeStub{}
var transactionContextInstance = &mocks.TransactionContext{}
var clientIdentityInstance = &mocks.ClientIdentity{MspID: "RFB-MINIFABRIC-LAB"}
var contract = chaincode.SmartContract{}

const txIDStub = "12345"

func TestMain(m *testing.M) {
	chaincodeStubInstance.GetTxIDReturns(txIDStub)
	transactionContextInstance.GetStubReturns(chaincodeStubInstance)
	clientIdentityInstance = &mocks.ClientIdentity{MspID: "RFB-MINIFABRIC-LAB"}
	transactionContextInstance.GetClientIdentityReturns(clientIdentityInstance)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestRegistra(t *testing.T) {
	txID, err := contract.Registra(transactionContextInstance, pessoafisicaJSON)

	require.NoError(t, err)

	var txIDJSON map[string]interface{}
	err = json.Unmarshal([]byte(txID), &txIDJSON)

	require.NoError(t, err)
	require.NotNil(t, txIDJSON["txID"].(string))
	require.Equal(t, txIDStub, txIDJSON["txID"].(string))
}

func TestRegistraErroApenasRFBRegistraPessoaFisica(t *testing.T) {
	clientIdentityInstance = &mocks.ClientIdentity{MspID: "DETRAN-MINIFABRIC-LAB"}
	transactionContextInstance.GetClientIdentityReturns(clientIdentityInstance)

	_, err := contract.Registra(transactionContextInstance, pessoafisicaJSON)

	require.Error(t, err)
	assert.Equal(t, fmt.Errorf(messages.TransacaoRestritaOrgao, "RFB"), err)
}

func TestRegistraErroJSONInvalido(t *testing.T) {
	clientIdentityInstance = &mocks.ClientIdentity{MspID: "RFB-MINIFABRIC-LAB"}
	transactionContextInstance.GetClientIdentityReturns(clientIdentityInstance)

	_, err := contract.Registra(transactionContextInstance, "a")

	require.Error(t, err)
	assert.Equal(t, fmt.Errorf(messages.JSONInvalido), err)
}

func TestRegistraErroAtributosFaltantes(t *testing.T) {
	_, err := contract.Registra(transactionContextInstance, "{}")

	require.Error(t, err)
}

func TestRegistraErroJSONAtributoInvalido(t *testing.T) {
	_, err := contract.Registra(transactionContextInstance, `{"pessoa": "fisica" }`)

	require.Error(t, err)
	assert.Equal(t, `json: unknown field "pessoa"`, err.Error())
}

func TestBusca(t *testing.T) {
	chaincodeStubInstance.GetStateReturns([]byte(pessoafisicaJSON), nil)

	pf, err := contract.Busca(transactionContextInstance, "12345")

	require.NoError(t, err)
	require.NotNil(t, pf)
}

func TestBuscaErroCpfInexistente(t *testing.T) {
	chaincodeStubInstance.GetStateReturns(nil, nil)

	_, err := contract.Busca(transactionContextInstance, "12345")

	require.Error(t, err)
	assert.Equal(t, fmt.Errorf(messages.NaoExistePessoafisicaCpf), err)
}

func TestBuscaErroInexperado(t *testing.T) {
	errorMock := fmt.Errorf("Fabric is offline!")
	chaincodeStubInstance.GetStateReturns(nil, errorMock)

	_, err := contract.Busca(transactionContextInstance, "12345")

	require.Error(t, err)
	assert.Equal(t, fmt.Errorf(messages.ErroBuscaPessoafisica, "12345", errorMock), err)
}

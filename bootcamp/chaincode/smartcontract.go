package chaincode

import (
	"encoding/json"
	"fmt"

	"bootcamp/chaincode/messages"
	"bootcamp/chaincode/utils"

	valid "github.com/asaskevich/govalidator"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract oferece transações para gerenciar uma Pessoa Fisica
type SmartContract struct {
	contractapi.Contract
}

// Pessoafisica descreve uma pessoa fisica registrada.
type Pessoafisica struct {
	Type           string `json:"type"`
	Cpf            string `json:"cpf" valid:"required~O atributo cpf é obrigatório"`
	NomeCompleto   string `json:"nomeCompleto" valid:"required~O atributo nomeCompleto é obrigatório"`
	DataNascimento string `json:"dataNascimento" valid:"required~O atributo dataNascimento é obrigatório"`
	Cnh            *Cnh   `json:"cnh,omitempty"`
}

// Cnh descreve uma cnh vinculada a uma Pessoafisica.
type Cnh struct {
}

// Registra cria uma nova pessoafisica (p) de acordo com string JSON passado por parâmetro.
func (c *SmartContract) Registra(ctx contractapi.TransactionContextInterface, p string) (string, error) {
	if err := utils.ValidateMSPID(ctx, "rfb"); err != nil {
		return "", err
	}

	pessoafisica := Pessoafisica{Type: "pessoafisica"}
	if err := utils.ParseJSON(p, &pessoafisica); err != nil {
		return "", err
	}
	if _, err := valid.ValidateStruct(pessoafisica); err != nil {
		return "", err
	}

	pessoafisicaJSON, _ := json.Marshal(pessoafisica)
	err := ctx.GetStub().PutState(pessoafisica.Cpf, pessoafisicaJSON)

	return fmt.Sprintf(`{"txID": "%s"}`, ctx.GetStub().GetTxID()), err
}

// Busca uma pessoafisica de cpf passado por parâmetro.
func (c *SmartContract) Busca(ctx contractapi.TransactionContextInterface, cpf string) (string, error) {
	pessoafisica, err := buscaPessoafisica(ctx, cpf)
	if err != nil {
		return "", err
	}

	pessoafisicaJSON, _ := json.Marshal(pessoafisica)

	return string(pessoafisicaJSON), nil
}

// VerificaMaioridade de uma pessoafisica de cpf passado por parâmetro.
func (c *SmartContract) VerificaMaioridade(ctx contractapi.TransactionContextInterface, cpf string) (string, error) {
	// 1. buscar a pessoa física;
	// 2. verificar idade. Dica: usar goment:
	//    dataNascimento, _ := goment.New(pessoafisica.DataNascimento)
	//    dataHoje.Diff(dataNascimento, "years");
	// 3. retornar true se for 'de maior', false caso contrário.
	return "", nil
}

// RegistraVinculaCNH inclui e vincula uma CNH a uma pessoa fisica de cpf passado por parâmetro.
func (c *SmartContract) RegistraVinculaCNH(ctx contractapi.TransactionContextInterface, cpf string, cnhJSON string) (string, error) {
	// 1. apenas o detran pode registrar e vincular CNH;
	// 2. efetuar o parse do JSON (cnhJSON) e validar;
	// 3. buscar a pessoa física;
	// 4. verificar se a pessoa física tem 18 anos ou mais;
	// 5. atribuir a Cnh passada por parâmetro a pessoa física encontrada;
	// 6. salvar a pessoa física;
	// 7. retornar txID.
	return fmt.Sprintf(`{"txID": "%s"}`, ctx.GetStub().GetTxID()), nil
}

// HealthFabric expõe uma transação para testar a conectividade.
func (c *SmartContract) HealthFabric(ctx contractapi.TransactionContextInterface) string {
	return `{ "status": "pass" }`
}

func buscaPessoafisica(ctx contractapi.TransactionContextInterface, cpf string) (*Pessoafisica, error) {
	pessoafisicaJSON, err := ctx.GetStub().GetState(cpf)
	if err != nil {
		return nil, fmt.Errorf(messages.ErroBuscaPessoafisica, cpf, err)
	}
	if pessoafisicaJSON == nil {
		return nil, fmt.Errorf(messages.NaoExistePessoafisicaCpf)
	}

	var pessoafisica Pessoafisica
	json.Unmarshal(pessoafisicaJSON, &pessoafisica)

	return &pessoafisica, nil
}

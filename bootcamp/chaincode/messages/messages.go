package messages

// JSONInvalido mensagem JSON inválido
const JSONInvalido string = "JSON inválido enviado para transação"

// TransacaoRestritaOrgao mensagem restrição de órgão
const TransacaoRestritaOrgao string = "Apenas %s pode efetuar essa transação!"

// ErroBuscaPessoafisica mensagem para erro inexperado na busca de pessoa física por CPF.
const ErroBuscaPessoafisica string = "Não foi possível carregar a pessoa fisica de CPF %s. Veja os erros: %v"

// NaoExistePessoafisicaCpf mensagem CPF inexistente na base quando é efetuada busca.
const NaoExistePessoafisicaCpf string = "Não existe pessoa física com o CPF fornecido na consulta!"

// MenorDeIdade mensagem emissão de CNH somente para pessoa física com idade >= 18.
const MenorDeIdade string = "CNH só pode ser emitida para maior de idade!"

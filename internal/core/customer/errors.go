package customer

import "fmt"

var (
	ErrMaxAddresses            = fmt.Errorf("você pode cadastrar até 5 endereços")
	ErrTryRemoveInvalidAddress = fmt.Errorf("o endereço não pertence a essa conta")
	ErrThereIsNoAddresses      = fmt.Errorf("você não possui este endereço cadastrado")
	ErrCustomerAlreadyExist    = fmt.Errorf("dados já cadastrados")
)

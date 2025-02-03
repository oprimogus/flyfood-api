package store

import "errors"

var (
	ErrClosingTimeBeforeOpeningTime = errors.New("o horário de encerramento não pode ser anterior ao horário de abertura no mesmo dia")
	ErrInvalidWeekDayBusinessHour   = errors.New("dia inválido")
	ErrBusinessHourAlreadyExist     = errors.New("horário de funcionamento já cadastrado")
	ErrBusinessHourNotExist         = errors.New("horário de funcionamento informado não existe")
	ErrRemoveInvalidPaymentMethod   = errors.New("essa loja  não possui o método de pagamento informado")
	ErrInvalidHour                  = errors.New("horário inválido")
	ErrOpeningHourEqualClosingHour  = errors.New("o horário de abertura e de encerramento não podem ser iguais")
	ErrPaymentMethodAlreadyDefined  = errors.New("essa loja já está habilitada com o método de pagamento informado")
)

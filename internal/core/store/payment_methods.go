package store

type PaymentMethod int

const (
	Credit  PaymentMethod = iota + 1
	Debit
	Cash
	Pix
	Bitcoin
)

func IsValidPaymentMethod(PaymentMethodEnum int) bool {
	switch PaymentMethod(PaymentMethodEnum) {
	case Credit, Debit, Pix, Cash, Bitcoin:
		return true
	default:
		return false
	}   
}

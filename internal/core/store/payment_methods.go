package store

type PaymentMethod string

const (
	Credit  PaymentMethod = "CREDIT"
	Debit   PaymentMethod = "DEBIT"
	Cash    PaymentMethod = "CASH"
	Pix     PaymentMethod = "PIX"
	Bitcoin PaymentMethod = "BTC"
)

func IsValidPaymentMethod(PaymentMethodEnum string) bool {
	switch PaymentMethod(PaymentMethodEnum) {
	case Credit, Debit, Pix, Cash, Bitcoin:
		return true
	default:
		return false
	}
}

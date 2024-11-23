package store

type Type string

const (
	Restaurant  Type = "RESTAURANT"
	Pharmacy    Type = "PHARMACY"
	Tobbaco     Type = "TOBBACO"
	Market      Type = "MARKET"
	Convenience Type = "CONVENIENCE"
	Pub         Type = "PUB"
)

func IsValidType(storeType string) bool {
	switch Type(storeType) {
	case Restaurant,
		Pharmacy,
		Tobbaco,
		Market,
		Convenience,
		Pub:
		return true
	default:
		return false
	}
}

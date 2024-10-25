package xvalidator

import (
	"reflect"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/oprimogus/cardapiogo/internal/core/item"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

func isValidUserRole(fl validator.FieldLevel) bool {
	role := fl.Field()

	if role.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < role.Len(); i++ {
		role := role.Index(i).String()
		if !user.IsValidRole(role) {
			return false
		}
	}

	return true
}

func IsValidCpfOrCnpj(fl validator.FieldLevel) bool {
	return IsValidCpf(fl) || IsValidCnpj(fl)
}

func IsValidPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	regex := `^\+(\d{2})(\d{2})(\d{9})$`
	re := regexp.MustCompile(regex)

	return re.MatchString(phone)
}

func IsValidShopType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return store.IsValidShopType(value)
}

func IsValidItemType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return item.IsValidItemType(value)
}

func IsValidPaymentMethod(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return store.IsValidPaymentMethod(value)
}

func IsValidBusinessHour(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return store.IsBusinessHourString(value)
}

func IsValidCpf(fl validator.FieldLevel) bool {
	cpf := fl.Field().String()

	if len(cpf) != 11 {
		return false
	}
	if isAllEqual(cpf) {
		return false
	}
	d1 := calculateDigitCpf(cpf, 10)
	d2 := calculateDigitCpf(cpf, 11)
	return strconv.Itoa(d1) == cpf[9:10] && strconv.Itoa(d2) == cpf[10:11]
}

func IsValidCnpj(fl validator.FieldLevel) bool {
	cnpj := fl.Field().String()

	if len(cnpj) != 14 {
		return false
	}
	if isAllEqual(cnpj) {
		return false
	}

	d1 := calculateDigitCnpj(cnpj, 12)
	d2 := calculateDigitCnpj(cnpj, 13)

	return strconv.Itoa(d1) == string(cnpj[12]) && strconv.Itoa(d2) == string(cnpj[13])
}

func isValidWeekDay(fl validator.FieldLevel) bool {
	weekDay := fl.Field().Int()
	if weekDay >= 0 && weekDay <= 6 {
		return true
	}
	return false
}

func isAllEqual(value string) bool {
	for i := range value {
		if value[i] != value[0] {
			return false
		}
	}
	return true
}

func calculateDigitCpf(cpf string, weight int) int {
	sum := 0
	count := weight - 1
	for i := 0; i < count; i++ {
		number, _ := strconv.Atoi(string(cpf[i]))
		sum += number * weight
		weight--
	}
	rest := sum % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}

func calculateDigitCnpj(cnpj string, factor int) int {
	sum := 0
	var weights []int
	if factor < 13 {
		weights = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	} else {
		weights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	}

	for i := 0; i < factor; i++ {
		num, _ := strconv.Atoi(string(cnpj[i]))
		sum += num * weights[i]
	}

	rest := sum % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}

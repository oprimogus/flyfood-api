package customer_test

// import (
// 	"testing"

// 	"github.com/oprimogus/flyfood-api/internal/core/address"
// 	"github.com/oprimogus/flyfood-api/internal/core/customer"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCustomer_SaveNewAddress(t *testing.T) {

// 	mockAddress := address.Address{
// 		Name:         "Endereço salvo",
// 		AddressLine1: "Rua 1",
// 		AddressLine2: "845 C3",
// 		Neighborhood: "Bairro teste",
// 		City:         "Guarujá",
// 		State:        "SP",
// 		PostalCode:   "11470-180",
// 		Country:      "BR",
// 	}

// 	c := &customer.Customer{
// 		ID:        "34876876358756",
// 		Name:      "John",
// 		LastName:  "Doe",
// 		CPF:       "",
// 		Email:     "johndoe@example.com",
// 		Phone:     "+5513997590579",
// 		Addresses: []address.Address{},
// 	}

// 	testcases := []struct {
// 		name          string
// 		customer      *customer.Customer
// 		actualAddress []address.Address
// 		newAddress    address.Address
// 		expected      error
// 	}{
// 		{
// 			name:     "Shoud add new address with success",
// 			customer: c,
// 			actualAddress: []address.Address{
// 				mockAddress,
// 				mockAddress,
// 			},
// 			newAddress: address.Address{
// 				Name:         "Casa 2",
// 				AddressLine1: "Rua 2",
// 				AddressLine2: "845 C3",
// 				Neighborhood: "Bairro teste",
// 				City:         "Guarujá",
// 				State:        "SP",
// 				PostalCode:   "11470-180",
// 				Country:      "BR",
// 			},
// 			expected: nil,
// 		},
// 		{
// 			name:     "Shoud return err when try add more than 5 address",
// 			customer: c,
// 			actualAddress: []address.Address{
// 				mockAddress,
// 				mockAddress,
// 				mockAddress,
// 				mockAddress,
// 				mockAddress,
// 			},
// 			newAddress: address.Address{
// 				Name:         "Trabalho",
// 				AddressLine1: "Rua 2",
// 				AddressLine2: "845 C3",
// 				Neighborhood: "Bairro teste",
// 				City:         "Guarujá",
// 				State:        "SP",
// 				PostalCode:   "11470-180",
// 				Country:      "BR",
// 			},
// 			expected: customer.ErrMaxAddresses,
// 		},
// 	}

// 	for _, test := range testcases {
// 		test.customer.Addresses = test.actualAddress
// 		err := test.customer.SaveNewAddress(test.newAddress)
// 		assert.Equal(t, test.expected, err)
// 	}
// }

// func TestCustomer_RemoveAddress(t *testing.T) {
// 	mockAddress := address.Address{
// 		Name:         "Endereço salvo",
// 		AddressLine1: "Rua 1",
// 		AddressLine2: "845 C3",
// 		Neighborhood: "Bairro teste",
// 		City:         "Guarujá",
// 		State:        "SP",
// 		PostalCode:   "11470-180",
// 		Country:      "BR",
// 	}

// 	customerWithAddress := &customer.Customer{
// 		ID:       "257547567254247245",
// 		Name:     "John",
// 		LastName: "Doe",
// 		CPF:      "",
// 		Email:    "johndoe@example.com",
// 		Phone:    "+5513997590579",
// 		Addresses: []address.Address{
// 			mockAddress,
// 		},
// 	}

// 	customerWithoutAddress := &customer.Customer{
// 		ID:        "2547724254724575472",
// 		Name:      "John",
// 		LastName:  "Doe",
// 		CPF:       "",
// 		Email:     "johndoe@example.com",
// 		Phone:     "+5513997590579",
// 		Addresses: []address.Address{},
// 	}

// 	testcases := []struct {
// 		name              string
// 		customer          *customer.Customer
// 		addressToRemove   address.Address
// 		expectedAddresses []address.Address
// 		expectedError     error
// 	}{
// 		{
// 			name:              "Should remove address",
// 			customer:          customerWithAddress,
// 			addressToRemove:   mockAddress,
// 			expectedAddresses: nil,
// 			expectedError:     nil,
// 		},
// 		{
// 			name:              "Should return error when try remove an address when customer no have any address",
// 			customer:          customerWithoutAddress,
// 			addressToRemove:   mockAddress,
// 			expectedAddresses: nil,
// 			expectedError:     customer.ErrThereIsNoAddresses,
// 		},
// 	}

// 	for _, test := range testcases {
// 		err := test.customer.RemoveAddress(test.addressToRemove)
// 		if test.expectedError == nil {
// 			assert.NoError(t, err)
// 			assert.Equal(t, test.expectedAddresses, test.customer.Addresses, test.name)
// 		} else {
// 			assert.Equal(t, test.expectedError, err, test.name)
// 		}
// 	}
// }

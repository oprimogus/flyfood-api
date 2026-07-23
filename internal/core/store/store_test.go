package store

import (
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/stretchr/testify/assert"
)

// createTestAddress é um helper para criar endereços de teste
func createTestAddress() address.Address {
	addressId, _ := uuid.NewV7()
	return address.Address{
		ID:           addressId,
		Name:         "Test Store Location",
		AddressLine1: "Test Street 123",
		AddressLine2: "Suite 100",
		Neighborhood: "Test Neighborhood",
		City:         "Test City",
		State:        "TS",
		PostalCode:   "12345-678",
		Country:      "Brasil",
	}
}

func TestNewStore(t *testing.T) {
	validOwnerID1, _ := uuid.NewV7()
	validOwnerID2, _ := uuid.NewV7()
	validOwnerID3, _ := uuid.NewV7()

	tests := []struct {
		name        string
		ownerID     uuid.UUID
		cnpj        string
		storeName   string
		description string
		phone       string
		address     address.Address
		storeType   Type
		expectError bool
	}{
		{
			name:        "should create new store successfully",
			ownerID:     validOwnerID1,
			cnpj:        "24611859000103",
			storeName:   "Test Restaurant",
			description: "A great test restaurant",
			phone:       "+5511997590670",
			address:     createTestAddress(),
			storeType:   Restaurant,
			expectError: false,
		},
		{
			name:        "should create new pharmacy store",
			ownerID:     validOwnerID2,
			cnpj:        "24611859000103",
			storeName:   "Test Pharmacy",
			description: "A reliable test pharmacy",
			phone:       "+5511997590670",
			address:     createTestAddress(),
			storeType:   Pharmacy,
			expectError: false,
		},
		{
			name:        "should create new market store",
			ownerID:     validOwnerID3,
			cnpj:        "24611859000103",
			storeName:   "Test Market",
			description: "A convenient test market",
			phone:       "+5511997590670",
			address:     createTestAddress(),
			storeType:   Market,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := NewStore(
				tt.ownerID,
				tt.cnpj,
				tt.storeName,
				tt.description,
				tt.phone,
				tt.address,
				tt.storeType,
			)

			slog.Error(tt.name, "error", err)

			if tt.expectError {
				slog.Error(tt.name, "error", err)
				assert.Error(t, err)
			}

			assert.NoError(t, err)
			assert.NotNil(t, store)

			// Verificar campos básicos
			assert.NotEmpty(t, store.ID, "Store ID should be generated")
			assert.Equal(t, tt.ownerID, store.OwnerID)
			assert.Equal(t, tt.cnpj, store.CNPJ)
			assert.Equal(t, tt.storeName, store.Name)
			assert.Equal(t, tt.description, store.Description)
			assert.Equal(t, tt.phone, store.Phone)
			assert.Equal(t, tt.storeType, store.Type)

			// Verificar valores padrão
			assert.Equal(t, DefaultScore, store.Score)
			assert.False(t, store.Active, "Store should default to inactive")
			assert.False(t, store.IsOpen, "Store should default to closed")
			assert.Empty(t, store.ProfileImage, "ProfileImage should be empty")
			assert.Empty(t, store.HeaderImage, "HeaderImage should be empty")
			assert.Empty(t, store.BusinessHours, "BusinessHours should be empty")
			assert.Empty(t, store.PaymentMethods, "PaymentMethods should be empty")
		})
	}
}

func TestStore_AddNewBusinessHour(t *testing.T) {
	a, _ := NewMinutesOfDayFromHHMM("13:00")
	b, _ := NewMinutesOfDayFromHHMM("18:00")
	c, _ := NewMinutesOfDayFromHHMM("12:00")
	d, _ := NewMinutesOfDayFromHHMM("21:00")

	tests := []struct {
		name               string
		actualBusinessHour []BusinessHours
		newBusinessHour    BusinessHours
		expected           error
	}{
		{
			name:               "should add new business hour",
			actualBusinessHour: []BusinessHours{},
			newBusinessHour: BusinessHours{
				WeekDay:     1,
				OpeningTime: a,
				ClosingTime: b,
			},
			expected: nil,
		},
		{
			name: "should return error when try add same business hour",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: c,
					ClosingTime: d,
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     1,
				OpeningTime: c,
				ClosingTime: d,
			},
			expected: ErrBusinessHourAlreadyExist,
		},
		{
			name: "should return error when try add invalid business hour",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: c,
					ClosingTime: d,
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     2,
				OpeningTime: a,
				ClosingTime: b,
			},
			expected: ErrClosingTimeBeforeOpeningTime,
		},
		{
			name: "should return error when try add invalid business hour where opening equal closing hour",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: c,
					ClosingTime: d,
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     2,
				OpeningTime: a,
				ClosingTime: a,
			},
			expected: ErrOpeningHourEqualClosingHour,
		},
		{
			name: "should add multiple different business hours",
			actualBusinessHour: []BusinessHours{
				{
					WeekDay:     1,
					OpeningTime: c,
					ClosingTime: d,
				},
			},
			newBusinessHour: BusinessHours{
				WeekDay:     2,
				OpeningTime: a,
				ClosingTime: b,
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, _ := uuid.NewV7()
			store, err := NewStore(
				id,
				"62722033000100",
				"Test Store",
				"Test Description",
				"+5511997590670",
				createTestAddress(),
				Restaurant,
			)
			slog.Error(test.name, slog.Any("err", err))
			assert.NoError(t, err)

			store.BusinessHours = test.actualBusinessHour

			err = store.AddNewBusinessHour(test.newBusinessHour)
			if test.expected != nil {
				assert.Error(t, test.expected, err)
			} else {
				assert.NoError(t, err)
			}

			if test.expected == nil {
				expectedCount := len(test.actualBusinessHour) + 1
				assert.Len(t, store.BusinessHours, expectedCount, "Business hour should be added")

				found := false
				for _, bh := range store.BusinessHours {
					if bh == test.newBusinessHour {
						found = true
						break
					}
				}
				assert.True(t, found, "New business hour should be present in the list")
			}
		})
	}
}

func TestStore_RemoveBusinessHour(t *testing.T) {

	a, _ := NewMinutesOfDayFromHHMM("09:00")
	b, _ := NewMinutesOfDayFromHHMM("17:00")

	tests := []struct {
		name          string
		initialHours  []BusinessHours
		hourToRemove  BusinessHours
		expectedError error
		expectedCount int
	}{
		{
			name: "should remove existing business hour",
			initialHours: []BusinessHours{
				{WeekDay: 1, OpeningTime: a, ClosingTime: b},
				{WeekDay: 2, OpeningTime: a, ClosingTime: b},
			},
			hourToRemove: BusinessHours{
				WeekDay: 1, OpeningTime: a, ClosingTime: b,
			},
			expectedError: nil,
			expectedCount: 1,
		},
		{
			name: "should return error when trying to remove non-existent business hour",
			initialHours: []BusinessHours{
				{WeekDay: 1, OpeningTime: a, ClosingTime: b},
			},
			hourToRemove: BusinessHours{
				WeekDay: 2, OpeningTime: a, ClosingTime: b,
			},
			expectedError: ErrBusinessHourNotExist,
			expectedCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, _ := uuid.NewV7()
			store, err := NewStore(
				id,
				"54560071000178",
				"Test Store",
				"Test Description",
				"+5511997590670",
				createTestAddress(),
				Restaurant,
			)
			assert.NoError(t, err)

			store.BusinessHours = test.initialHours

			err = store.RemoveBusinessHour(test.hourToRemove)
			assert.Equal(t, test.expectedError, err)
			assert.Len(t, store.BusinessHours, test.expectedCount)
		})
	}
}

func TestStore_AddPaymentMethod(t *testing.T) {
	tests := []struct {
		name                  string
		initialPaymentMethods []PaymentMethod
		paymentMethodToAdd    PaymentMethod
		expectedError         error
		expectedCount         int
	}{
		{
			name:                  "should add new payment method",
			initialPaymentMethods: []PaymentMethod{},
			paymentMethodToAdd:    Credit,
			expectedError:         nil,
			expectedCount:         1,
		},
		{
			name:                  "should add multiple payment methods",
			initialPaymentMethods: []PaymentMethod{Credit},
			paymentMethodToAdd:    Pix,
			expectedError:         nil,
			expectedCount:         2,
		},
		{
			name:                  "should return error when adding duplicate payment method",
			initialPaymentMethods: []PaymentMethod{Credit, Pix},
			paymentMethodToAdd:    Credit,
			expectedError:         ErrPaymentMethodAlreadyDefined,
			expectedCount:         2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, _ := uuid.NewV7()
			store, err := NewStore(
				id,
				"34021317000126",
				"Test Store",
				"Test Description",
				"+5511997590670",
				createTestAddress(),
				Restaurant,
			)
			assert.NoError(t, err)

			store.PaymentMethods = test.initialPaymentMethods

			err = store.AddPaymentMethod(test.paymentMethodToAdd)
			assert.Equal(t, test.expectedError, err)
			assert.Len(t, store.PaymentMethods, test.expectedCount)

			if test.expectedError == nil {
				found := false
				for _, pm := range store.PaymentMethods {
					if pm == test.paymentMethodToAdd {
						found = true
						break
					}
				}
				assert.True(t, found, "Payment method should be added")
			}
		})
	}
}

func TestStore_RemovePaymentMethod(t *testing.T) {
	tests := []struct {
		name                  string
		initialPaymentMethods []PaymentMethod
		paymentMethodToRemove PaymentMethod
		expectedError         error
		expectedCount         int
	}{
		{
			name:                  "should remove existing payment method",
			initialPaymentMethods: []PaymentMethod{Credit, Pix, Cash},
			paymentMethodToRemove: Pix,
			expectedError:         nil,
			expectedCount:         2,
		},
		{
			name:                  "should return error when removing non-existent payment method",
			initialPaymentMethods: []PaymentMethod{Credit},
			paymentMethodToRemove: Bitcoin,
			expectedError:         ErrRemoveInvalidPaymentMethod,
			expectedCount:         1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, _ := uuid.NewV7()
			store, err := NewStore(
				id,
				"41614248000120",
				"Test Store",
				"Test Description",
				"+5511997590670",
				createTestAddress(),
				Restaurant,
			)
			assert.NoError(t, err)

			store.PaymentMethods = test.initialPaymentMethods

			err = store.RemovePaymentMethod(test.paymentMethodToRemove)
			assert.Equal(t, test.expectedError, err)
			assert.Len(t, store.PaymentMethods, test.expectedCount)
		})
	}
}

func TestStore_StateChanges(t *testing.T) {
	t.Run("should activate and deactivate store", func(t *testing.T) {
		id, _ := uuid.NewV7()
		store, err := NewStore(
			id,
			"12929432000160",
			"Test Store",
			"Test Description",
			"+5511997590670",
			createTestAddress(),
			Restaurant,
		)
		assert.NoError(t, err)
		assert.False(t, store.Active, "Store should start inactive")

		store.Activate()
		assert.True(t, store.Active, "Store should be active after activation")

		store.Deactivate()
		assert.False(t, store.Active, "Store should be inactive after deactivation")
	})

	t.Run("should open and close store", func(t *testing.T) {
		id, _ := uuid.NewV7()
		store, err := NewStore(
			id,
			"47416049000193",
			"Test Store",
			"Test Description",
			"+5511997590670",
			createTestAddress(),
			Restaurant,
		)
		assert.NoError(t, err)
		assert.False(t, store.IsOpen, "Store should start closed")

		store.OpenStore()
		assert.True(t, store.IsOpen, "Store should be open after opening")

		store.CloseStore()
		assert.False(t, store.IsOpen, "Store should be closed after closing")
	})
}

func TestStore_UpdateStore(t *testing.T) { // Corrigido o nome do teste e do método chamado
	id, _ := uuid.NewV7()
	addressId, _ := uuid.NewV7()
	store, err := NewStore(
		id,
		"08151022000164",
		"Original Name",
		"Original Description",
		"+5511997590670",
		createTestAddress(),
		Restaurant,
	)
	assert.NoError(t, err)

	newAddress := address.Address{
		ID:           addressId,
		Name:         "Updated Location",
		AddressLine1: "Updated Street 456",
		AddressLine2: "Updated Suite 200",
		Neighborhood: "Updated Neighborhood",
		City:         "Updated City",
		State:        "UP",
		PostalCode:   "98765-432",
		Country:      "Brasil",
	}

	// Método correto é UpdateStore e não aceita número (delivery time) no final
	err = store.UpdateStore(
		"Updated Store Name",
		"Updated Description",
		"+5511997590670",
		newAddress,
		Pharmacy,
	)
	slog.Error("TestStore_UpdateStore", slog.Any("err", err))
	assert.NoError(t, err)

	assert.Equal(t, "Updated Store Name", store.Name)
	assert.Equal(t, "Updated Description", store.Description)
	assert.Equal(t, "+5511997590670", store.Phone)
	assert.Equal(t, newAddress, store.Address)
	assert.Equal(t, Pharmacy, store.Type)
}

func TestStore_ImageMethods(t *testing.T) {
	id, _ := uuid.NewV7()
	store, err := NewStore(
		id,
		"26408539000178",
		"Test Store",
		"Test Description",
		"+5511997590670",
		createTestAddress(),
		Restaurant,
	)
	assert.NoError(t, err)

	t.Run("should change profile image", func(t *testing.T) {
		profileImageURL := "https://example.com/profile.jpg"
		store.ChangeProfileImage(profileImageURL)
		assert.Equal(t, profileImageURL, store.ProfileImage)
	})

	t.Run("should change header image", func(t *testing.T) {
		headerImageURL := "https://example.com/header.jpg"
		store.ChangeHeaderImage(headerImageURL)
		assert.Equal(t, headerImageURL, store.HeaderImage)
	})
}

func TestStore_AllStoreTypes(t *testing.T) {
	storeTypes := []Type{Restaurant, Pharmacy, Tobacco, Market, Convenience, Pub}

	for _, storeType := range storeTypes {
		id, _ := uuid.NewV7()
		t.Run(storeType.String(), func(t *testing.T) {
			store, err := NewStore(
				id,
				"83157985000190",
				"Test "+storeType.String(),
				"Test Description",
				"+5511997590670",
				createTestAddress(),
				storeType,
			)
			assert.NoError(t, err)
			assert.Equal(t, storeType, store.Type)
		})
	}
}

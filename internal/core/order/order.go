package order

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/oprimogus/flyfood-api/internal/core/address"
	"github.com/oprimogus/flyfood-api/pkg/validator"
)

type State string

const (
	Created            State = "created"              // Order created by customer
	Cancelled          State = "cancelled"            // Order cancelled by customer or store
	VerifiedByCustomer State = "verified_by_customer" // Order confirmed by the customer that the order is correct
	VerifiedByStore    State = "verified_by_store"    // Order confirmed by the Store and will be processed
	InProcess          State = "in_process"           // Order in process
	Dispatched         State = "dispatched"           // Order dispatched for Order address
	Delivered          State = "delivered"            // Order was delivered to customer
	ChargeBack         State = "charge_back"          // Customer want chargeback of order
	Finished           State = "finished"             // Order finished
)

type Order struct {
	ID             string          `json:"id" validate:"required,uuid"`
	StoreID        string          `json:"store_id" validate:"required,uuid"`
	CustomerID     string          `json:"customer_id" validate:"required,uuid"`
	CourierID      string          `json:"courier_id" validate:"uuid"`
	Status         State           `json:"status" validate:"required,orderStatus"`
	Items          []Item          `json:"items" validate:"required,dive"`
	Amount         int             `json:"amount" validate:"required,gt=0"`
	DeliveryAmount int             `json:"deliveryAmount" validate:"required,gte=0"`
	Address        address.Address `json:"address" validate:"required,address"`
}

func (o *Order) Validate() error {
	return validator.Validate(o)
}

func NewOrder(storeID, customerID string, items []Item, amount, deliveryAmount int, address address.Address) (*Order, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("não foi possível gerar ID único para essa ordem: %w", err)
	}

	o := Order{
		ID:             id.String(),
		StoreID:        storeID,
		CustomerID:     customerID,
		Status:         Created,
		Items:          items,
		Amount:         amount,
		DeliveryAmount: deliveryAmount,
		Address:        address,
	}

	if err := o.Validate(); err != nil {
		return nil, err
	}

	return &o, nil
}

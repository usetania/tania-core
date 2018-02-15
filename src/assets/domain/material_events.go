package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MaterialCreated struct {
	UID            uuid.UUID
	Name           string
	PricePerUnit   Money
	Type           MaterialType
	Quantity       MaterialQuantity
	ExpirationDate *time.Time
	Notes          *string
	IsExpense      *bool
	ProducedBy     *string
}

type MaterialNameChanged struct {
	UID  uuid.UUID
	Name string
}

type MaterialPriceChanged struct {
	UID   uuid.UUID
	Price Money
}

type MaterialQuantityChanged struct {
	UID      uuid.UUID
	Quantity MaterialQuantity
}

type MaterialExpirationDateChanged struct {
	UID            uuid.UUID
	ExpirationDate time.Time
}

type MaterialNotesChanged struct {
	UID   uuid.UUID
	Notes string
}

type MaterialProducedByChanged struct {
	UID        uuid.UUID
	ProducedBy string
}

type MaterialIsExpenseChanged struct {
	UID       uuid.UUID
	IsExpense bool
}

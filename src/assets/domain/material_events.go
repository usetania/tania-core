package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type MaterialCreated struct {
	UID            uuid.UUID
	Name           string
	PricePerUnit   PricePerUnit
	Type           MaterialType
	Quantity       MaterialQuantity
	ExpirationDate *time.Time
	Notes          *string
	ProducedBy     *string
	CreatedDate    time.Time
}

type MaterialNameChanged struct {
	MaterialUID uuid.UUID
	Name        string
}

type MaterialPriceChanged struct {
	MaterialUID uuid.UUID
	Price       PricePerUnit
}

type MaterialQuantityChanged struct {
	MaterialUID      uuid.UUID
	MaterialTypeCode string
	Quantity         MaterialQuantity
}

type MaterialTypeChanged struct {
	MaterialUID  uuid.UUID
	MaterialType MaterialType
}

type MaterialExpirationDateChanged struct {
	MaterialUID    uuid.UUID
	ExpirationDate time.Time
}

type MaterialNotesChanged struct {
	MaterialUID uuid.UUID
	Notes       string
}

type MaterialProducedByChanged struct {
	MaterialUID uuid.UUID
	ProducedBy  string
}

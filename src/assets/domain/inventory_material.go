package domain

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Material struct {
	UID            uuid.UUID
	Name           string
	PricePerUnit   Money
	Type           MaterialType
	Quantity       float32
	QuantityUnit   MaterialQuantityUnit
	ExpirationDate time.Time
	Notes          string
	IsExpense      bool
	ProducedBy     string
}

const (
	MoneyEUR = "EUR"
	MoneyIDR = "IDR"
)

type Money interface {
	Code() string
	Symbol() string
	Amount() string
	SetAmount(amount string)
}

type EUR struct {
	amount string
}

func (e EUR) Code() string {
	return MoneyEUR
}

func (e EUR) Symbol() string {
	return "€"
}

func (e EUR) Amount() string {
	return e.amount
}

func (e *EUR) SetAmount(amount string) {
	e.amount = amount
}

func CreateMoney(price, priceUnit string) (Money, error) {
	if price == "" {
		return nil, errors.New("price cannot be empty")
	}

	var money Money
	switch priceUnit {
	case EUR{}.Code():
		money = &EUR{}
		money.SetAmount(price)
	default:
		return nil, errors.New("money not found")
	}

	return money, nil
}

const (
	MaterialUnitSeeds    = "SEEDS"
	MaterialUnitPackets  = "PACKETS"
	MaterialUnitGram     = "GRAM"
	MaterialUnitKilogram = "KILOGRAM"
)

type MaterialQuantityUnit struct {
	Code  string
	Label string
}

func MaterialTypeSeedQuantityUnits() []MaterialQuantityUnit {
	return []MaterialQuantityUnit{
		{Code: MaterialUnitSeeds, Label: "Seeds"},
		{Code: MaterialUnitPackets, Label: "Packets"},
		{Code: MaterialUnitGram, Label: "Gram"},
		{Code: MaterialUnitKilogram, Label: "Kilogram"},
	}
}

func GetMaterialTypeSeedQuantityUnit(code string) MaterialQuantityUnit {
	for _, v := range MaterialTypeSeedQuantityUnits() {
		if v.Code == code {
			return v
		}
	}

	return MaterialQuantityUnit{}
}

func GetQuantityUnit(quantityUnit string, materialType MaterialType) MaterialQuantityUnit {
	switch materialType.(type) {
	case MaterialTypeSeed:
		return GetMaterialTypeSeedQuantityUnit(quantityUnit)
	}

	return MaterialQuantityUnit{}
}

func CreateMaterial(
	name string,
	price string,
	priceUnit string,
	materialType MaterialType,
	quantity float32,
	quantityUnit string) (Material, error) {

	uid, err := uuid.NewV4()
	if err != nil {
		return Material{}, err
	}

	money, err := CreateMoney(price, priceUnit)
	if err != nil {
		return Material{}, err
	}

	if materialType == nil {
		return Material{}, errors.New("cannot be empty")
	}

	qu := GetQuantityUnit(quantityUnit, materialType)
	if qu == (MaterialQuantityUnit{}) {
		return Material{}, errors.New("Cannot be empty")
	}

	return Material{
		UID:          uid,
		Name:         name,
		PricePerUnit: money,
		Type:         materialType,
		Quantity:     quantity,
		QuantityUnit: qu,
	}, nil
}

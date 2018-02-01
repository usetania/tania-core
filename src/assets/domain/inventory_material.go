package domain

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Material struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	PricePerUnit   Money            `json:"price_per_unit"`
	Type           MaterialType     `json:"type"`
	Quantity       MaterialQuantity `json:"quantity"`
	ExpirationDate *time.Time       `json:"expiration_date"`
	Notes          *string          `json:"notes"`
	IsExpense      *bool            `json:"is_expense"`
	ProducedBy     *string          `json:"produced_by"`
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
	MaterialUnitSeeds      = "SEEDS"
	MaterialUnitPackets    = "PACKETS"
	MaterialUnitGram       = "GRAM"
	MaterialUnitKilogram   = "KILOGRAM"
	MaterialUnitBags       = "BAGS"
	MaterialUnitBottles    = "BOTTLES"
	MaterialUnitCubicMetre = "CUBIC_METRE"
)

type MaterialQuantity struct {
	Value float32              `json:"value"`
	Unit  MaterialQuantityUnit `json:"unit"`
}

type MaterialQuantityUnit struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

func MaterialQuantityUnits(materialTypeCode string) []MaterialQuantityUnit {
	switch materialTypeCode {
	case MaterialTypeSeedCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitSeeds, Label: "Seeds"},
			{Code: MaterialUnitPackets, Label: "Packets"},
			{Code: MaterialUnitGram, Label: "Gram"},
			{Code: MaterialUnitKilogram, Label: "Kilogram"},
		}
	case MaterialTypeAgrochemicalCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitPackets, Label: "Packets"},
			{Code: MaterialUnitBottles, Label: "Bottles"},
			{Code: MaterialUnitBags, Label: "Bags"},
		}
	case MaterialTypeGrowingMediumCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitBags, Label: "Bags"},
			{Code: MaterialUnitCubicMetre, Label: "Cubic Metre"},
		}
	}

	return nil
}

func GetMaterialQuantityUnit(materialTypeCode string, code string) MaterialQuantityUnit {
	for _, v := range MaterialQuantityUnits(materialTypeCode) {
		if v.Code == code {
			return v
		}
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

	err = validateQuantity(quantity)
	if err != nil {
		return Material{}, err
	}

	qu, err := validateQuantityUnit(quantityUnit, materialType)
	if err != nil {
		return Material{}, err
	}

	return Material{
		UID:          uid,
		Name:         name,
		PricePerUnit: money,
		Type:         materialType,
		Quantity: MaterialQuantity{
			Value: quantity,
			Unit:  qu,
		},
	}, nil
}

func (m *Material) ChangeName(name string) error {
	if name == "" {
		return errors.New("cannot be empty")
	}

	if len(name) <= 5 {
		return errors.New("too few characters")
	}

	m.Name = name

	return nil
}

func (m *Material) ChangePricePerUnit(price, priceUnit string) error {
	money, err := CreateMoney(price, priceUnit)
	if err != nil {
		return err
	}

	m.PricePerUnit = money

	return nil
}

func (m *Material) ChangeQuantityUnit(quantity float32, quantityUnit string, materialType MaterialType) error {
	err := validateQuantity(quantity)
	if err != nil {
		return err
	}

	qu, err := validateQuantityUnit(quantityUnit, materialType)
	if err != nil {
		return err
	}

	m.Quantity = MaterialQuantity{
		Value: quantity,
		Unit:  qu,
	}

	return nil
}

func validateQuantity(quantity float32) error {
	if quantity <= 0 {
		return errors.New("Cannot be empty")
	}

	return nil
}

func validateQuantityUnit(quantityUnit string, materialType MaterialType) (MaterialQuantityUnit, error) {
	qu := GetMaterialQuantityUnit(materialType.Code(), quantityUnit)

	if qu == (MaterialQuantityUnit{}) {
		return MaterialQuantityUnit{}, errors.New("Cannot be empty")
	}

	return qu, nil
}

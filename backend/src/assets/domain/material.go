package domain

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type Material struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	PricePerUnit   PricePerUnit     `json:"price_per_unit"`
	Type           MaterialType     `json:"type"`
	Quantity       MaterialQuantity `json:"quantity"`
	ExpirationDate *time.Time       `json:"expiration_date"`
	Notes          *string          `json:"notes"`
	ProducedBy     *string          `json:"produced_by"`
	CreatedDate    time.Time        `json:"created_date"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

const (
	MoneyEUR = "EUR"
	MoneyIDR = "IDR"
)

type PricePerUnit struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"code"`
}

func (p PricePerUnit) Symbol() string {
	switch p.CurrencyCode {
	case MoneyEUR:
		return "€"
	default:
		return ""
	}
}

func CreatePricePerUnit(amount, currencyCode string) (PricePerUnit, error) {
	cc, err := GetCurrencyCode(currencyCode)
	if err != nil {
		return PricePerUnit{}, err
	}

	return PricePerUnit{
		Amount:       amount,
		CurrencyCode: cc,
	}, nil
}

func GetCurrencyCode(currencyCode string) (string, error) {
	switch currencyCode {
	case MoneyEUR:
		return MoneyEUR, nil
	default:
		return "", errors.New("wrong currency code")
	}
}

const (
	MaterialUnitSeeds      = "SEEDS"
	MaterialUnitPackets    = "PACKETS"
	MaterialUnitGram       = "GRAM"
	MaterialUnitKilogram   = "KILOGRAM"
	MaterialUnitBags       = "BAGS"
	MaterialUnitBottles    = "BOTTLES"
	MaterialUnitCubicMetre = "CUBIC_METRE"
	MaterialUnitPieces     = "PIECES"
	MaterialUnitUnits      = "UNITS"
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
	case MaterialTypeLabelAndCropSupportCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitPieces, Label: "Pieces"},
		}
	case MaterialTypeSeedingContainerCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitPieces, Label: "Pieces"},
		}
	case MaterialTypePostHarvestSupplyCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitPieces, Label: "Pieces"},
		}
	case MaterialTypePlantCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitUnits, Label: "Units"},
			{Code: MaterialUnitPackets, Label: "Packets"},
		}
	case MaterialTypeOtherCode:
		return []MaterialQuantityUnit{
			{Code: MaterialUnitPieces, Label: "Pieces"},
		}
	}

	return nil
}

func GetMaterialQuantityUnit(materialTypeCode, code string) MaterialQuantityUnit {
	for _, v := range MaterialQuantityUnits(materialTypeCode) {
		if v.Code == code {
			return v
		}
	}

	return MaterialQuantityUnit{}
}

func (m *Material) TrackChange(event interface{}) {
	m.UncommittedChanges = append(m.UncommittedChanges, event)
	m.Transition(event)
}

func (m *Material) Transition(event interface{}) {
	switch e := event.(type) {
	case MaterialCreated:
		m.UID = e.UID
		m.Name = e.Name
		m.PricePerUnit = e.PricePerUnit
		m.Type = e.Type
		m.Quantity = e.Quantity
		m.ExpirationDate = e.ExpirationDate
		m.Notes = e.Notes
		m.ProducedBy = e.ProducedBy
		m.CreatedDate = e.CreatedDate

	case MaterialNameChanged:
		m.Name = e.Name

	case MaterialTypeChanged:
		m.Type = e.MaterialType

	case MaterialPriceChanged:
		m.PricePerUnit = e.Price

	case MaterialQuantityChanged:
		m.Quantity = e.Quantity

	case MaterialExpirationDateChanged:
		m.ExpirationDate = &e.ExpirationDate

	case MaterialNotesChanged:
		m.Notes = &e.Notes

	case MaterialProducedByChanged:
		m.ProducedBy = &e.ProducedBy
	}
}

func CreateMaterial(
	name string,
	price string,
	priceUnit string,
	materialType MaterialType,
	quantity float32,
	quantityUnit string,
	expirationDate *time.Time,
	notes *string,
	producedBy *string) (*Material, error,
) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	pricePerUnit, err := CreatePricePerUnit(price, priceUnit)
	if err != nil {
		return nil, err
	}

	if materialType == nil {
		return nil, errors.New("cannot be empty")
	}

	err = validateQuantity(quantity)
	if err != nil {
		return nil, err
	}

	qu, err := validateQuantityUnit(quantityUnit, materialType)
	if err != nil {
		return nil, err
	}

	initial := &Material{
		UID:          uid,
		Name:         name,
		PricePerUnit: pricePerUnit,
		Type:         materialType,
		Quantity: MaterialQuantity{
			Value: quantity,
			Unit:  qu,
		},
		ExpirationDate: expirationDate,
		Notes:          notes,
		ProducedBy:     producedBy,
		CreatedDate:    time.Now(),
	}

	initial.TrackChange(MaterialCreated{
		UID:            initial.UID,
		Name:           initial.Name,
		PricePerUnit:   initial.PricePerUnit,
		Type:           initial.Type,
		Quantity:       initial.Quantity,
		ExpirationDate: initial.ExpirationDate,
		Notes:          initial.Notes,
		ProducedBy:     initial.ProducedBy,
		CreatedDate:    initial.CreatedDate,
	})

	return initial, nil
}

func (m *Material) ChangeName(name string) error {
	if name == "" {
		return errors.New("cannot be empty")
	}

	if len(name) <= 5 {
		return errors.New("too few characters")
	}

	m.TrackChange(MaterialNameChanged{MaterialUID: m.UID, Name: name})

	return nil
}

func (m *Material) ChangePricePerUnit(price, priceUnit string) error {
	ppu, err := CreatePricePerUnit(price, priceUnit)
	if err != nil {
		return err
	}

	m.TrackChange(MaterialPriceChanged{MaterialUID: m.UID, Price: ppu})

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

	m.TrackChange(MaterialQuantityChanged{
		MaterialUID: m.UID,
		Quantity: MaterialQuantity{
			Value: quantity,
			Unit:  qu,
		},
		MaterialTypeCode: materialType.Code(),
	})

	return nil
}

func (m *Material) ChangeType(materialType MaterialType) error {
	if materialType == nil {
		return MaterialError{MaterialErrorInvalidMaterialType}
	}

	m.TrackChange(MaterialTypeChanged{
		MaterialUID:  m.UID,
		MaterialType: materialType,
	})

	return nil
}

func (m *Material) ChangeExpirationDate(expDate time.Time) error {
	m.TrackChange(MaterialExpirationDateChanged{
		MaterialUID:    m.UID,
		ExpirationDate: expDate,
	})

	return nil
}

func (m *Material) ChangeNotes(notes string) error {
	m.TrackChange(MaterialNotesChanged{
		MaterialUID: m.UID,
		Notes:       notes,
	})

	return nil
}

func (m *Material) ChangeProducedBy(producedBy string) error {
	m.TrackChange(MaterialProducedByChanged{
		MaterialUID: m.UID,
		ProducedBy:  producedBy,
	})

	return nil
}

func validateQuantity(quantity float32) error {
	if quantity <= 0 {
		return errors.New("cannot be empty")
	}

	return nil
}

func validateQuantityUnit(quantityUnit string, materialType MaterialType) (MaterialQuantityUnit, error) {
	qu := GetMaterialQuantityUnit(materialType.Code(), quantityUnit)

	if qu == (MaterialQuantityUnit{}) {
		return MaterialQuantityUnit{}, errors.New("cannot be empty")
	}

	return qu, nil
}

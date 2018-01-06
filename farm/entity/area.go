package entity

import (
	"github.com/Tanibox/tania-server/helper/validationhelper"
	uuid "github.com/satori/go.uuid"
)

type Area struct {
	UID       uuid.UUID `json:"uid"`
	Name      string    `json:"name"`
	Size      AreaUnit  `json:"size"`
	Type      string    `json:"type"`
	Location  string    `json:"location"`
	Photo     AreaPhoto `json:"photo"`
	Reservoir Reservoir `json:"reservoir"`
	Farm      Farm      `json:"-"`
}

type AreaUnit interface {
	Symbol() string
}

type SquareMeter struct {
	Value float32
}

func (sm SquareMeter) Symbol() string {
	return "m2"
}

type Hectare struct {
	Value float32
}

func (h Hectare) Symbol() string {
	return "Ha"
}

type AreaPhoto struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// CreateArea registers a new area to a farm
func CreateArea(farm Farm, name string, areaType string) (Area, error) {
	err := validateAreaName(name)
	if err != nil {
		return Area{}, err
	}

	err = validateAreaType(areaType)
	if err != nil {
		return Area{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Area{}, err
	}

	return Area{
		UID:  uid,
		Farm: farm,
		Name: name,
		Type: areaType,
	}, nil
}

// ChangeSize changes an area size
func (a *Area) ChangeSize(size AreaUnit) error {
	err := validateSize(size)
	if err != nil {
		return err
	}

	a.Size = size

	return nil
}

// ChangeLocation changes an area location
func (a *Area) ChangeLocation(location string) error {
	_, err := FindAreaLocationByCode(location)
	if err != nil {
		return err
	}

	a.Location = location

	return nil
}

func validateAreaName(name string) error {
	if name == "" {
		return AreaError{AreaErrorNameEmptyCode}
	}
	if !validationhelper.IsAlphanumeric(name) {
		return AreaError{AreaErrorNameAlphanumericOnlyCode}
	}
	if len(name) < 5 {
		return AreaError{AreaErrorNameNotEnoughCharacterCode}
	}
	if len(name) > 100 {
		return AreaError{AreaErrorNameExceedMaximunCharacterCode}
	}

	return nil
}

func validateSize(size AreaUnit) error {
	isValidUnit := true
	sizeValue := float32(0)

	switch v := size.(type) {
	case SquareMeter:
		sizeValue = v.Value
		isValidUnit = true
	case Hectare:
		sizeValue = v.Value
		isValidUnit = true
	default:
		sizeValue = 0
		isValidUnit = false
	}

	if sizeValue <= 0 {
		return AreaError{AreaErrorSizeEmptyCode}
	}
	if isValidUnit == false {
		return AreaError{AreaErrorInvalidSizeUnitCode}
	}

	return nil
}

func validateAreaType(areaType string) error {
	if areaType == "" {
		return AreaError{AreaErrorTypeEmptyCode}
	}

	_, err := FindAreaTypeByCode(areaType)
	if err != nil {
		return err
	}

	return nil
}

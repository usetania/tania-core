package entity

import "github.com/Tanibox/tania-server/helper/validationhelper"

type Area struct {
	UID       string      `json:"uid"`
	Name      string      `json:"name"`
	Size      float32     `json:"size"`
	SizeUnit  interface{} `json:"size_unit"`
	Type      string      `json:"type"`
	Location  string      `json:"location"`
	Photo     AreaPhoto   `json:"photo"`
	Reservoir Reservoir   `json:"reservoir"`
	Farm      Farm        `json:"-"`
}

type AreaSizeUnitMetre struct{}
type AreaSizeUnitHectare struct{}

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

	return Area{
		Farm: farm,
		Name: name,
		Type: areaType,
	}, nil
}

// ChangeSize changes an area size
func (a *Area) ChangeSize(size float32) error {
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

func validateSize(size float32) error {
	if size <= 0 {
		return AreaError{AreaErrorSizeEmptyCode}
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

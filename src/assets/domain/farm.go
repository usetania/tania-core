// Package domain provides the operation that farm holder can do
// to their farm
package domain

import (
	uuid "github.com/satori/go.uuid"
)

type Farm struct {
	UID         uuid.UUID `json:"uid"`
	Name        string    `json:"name"`
	Latitude    string    `json:"latitude"`
	Longitude   string    `json:"longitude"`
	Type        string    `json:"type"`
	CountryCode string    `json:"country_code"`
	CityCode    string    `json:"city_code"`
	IsActive    bool      `json:"is_active"`

	Reservoirs []Reservoir `json:"-"`
	Areas      []Area      `json:"-"`
}

// CreateFarm registers a new farm to Tania
func CreateFarm(name string, farmType string) (Farm, error) {
	err := validateFarmName(name)
	if err != nil {
		return Farm{}, err
	}

	err = validateFarmType(farmType)
	if err != nil {
		return Farm{}, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return Farm{}, err
	}

	return Farm{
		UID:         uid,
		Name:        name,
		Type:        farmType,
		Latitude:    "",
		Longitude:   "",
		CountryCode: "",
		CityCode:    "",
		IsActive:    true,
	}, nil
}

// ChangeGeoLocation changes the geolocation of a farm
func (f *Farm) ChangeGeoLocation(latitude, longitude string) error {
	err := validateGeoLocation(latitude, longitude)
	if err != nil {
		return err
	}

	f.Latitude = latitude
	f.Longitude = longitude

	return nil
}

// ChangeRegion changes country and city of a farm
func (f *Farm) ChangeRegion(countryCode, cityCode string) error {
	err := validateCountryCode(countryCode)
	if err != nil {
		return err
	}

	err = validateCityCode(countryCode, cityCode)
	if err != nil {
		return err
	}

	f.CountryCode = countryCode
	f.CityCode = cityCode

	return nil
}

// AddReservoir adds a reservoir to a farm
func (f *Farm) AddReservoir(res *Reservoir) error {
	if f.IsReservoirAdded(res.Name) {
		return FarmError{FarmErrorReservoirAlreadyAdded}
	}

	f.Reservoirs = append(f.Reservoirs, *res)

	return nil
}

// ChangeReservoirInformation changes existing reservoir information
func (f *Farm) ChangeReservoirInformation(res Reservoir) error {
	if !f.IsReservoirAdded(res.Name) {
		return FarmError{FarmErrorReservoirNotFound}
	}

	for i, v := range f.Reservoirs {
		if v.UID == res.UID {
			f.Reservoirs[i] = res
		}
	}

	return nil
}

// IsReservoirAdded is to check whether a reservoir is already added.
// It knows by matching the reservoir's name
func (f Farm) IsReservoirAdded(name string) bool {
	for _, r := range f.Reservoirs {
		if r.Name == name {
			return true
		}
	}

	return false
}

// IsHaveReservoir checks whether a farm has any reservoir.
func (f Farm) IsHaveReservoir() bool {
	if len(f.Reservoirs) > 0 {
		return true
	}

	return false
}

// AddArea adds a area to a farm
func (f *Farm) AddArea(res *Area) error {
	if f.IsAreaAdded(res.Name) {
		return FarmError{FarmErrorAreaAlreadyAdded}
	}

	f.Areas = append(f.Areas, *res)

	return nil
}

// IsAreaAdded is to check whether a area is already added.
// It knows by matching the area's name
func (f Farm) IsAreaAdded(name string) bool {
	for _, r := range f.Areas {
		if r.Name == name {
			return true
		}
	}

	return false
}

// IsHaveArea checks whether a farm has any area.
func (f Farm) IsHaveArea() bool {
	if len(f.Areas) > 0 {
		return true
	}

	return false
}

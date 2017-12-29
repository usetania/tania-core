// Package entity provides the operation that farm holder can do
// to their farm
package entity

type Farm struct {
	UID         string
	Name        string
	Latitude    string
	Longitude   string
	Type        string
	CountryCode string
	CityCode    string
	IsActive    bool

	Reservoirs []Reservoir
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

	return Farm{
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

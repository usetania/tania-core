// Package entity provides the operation that farm holder can do
// to their farm
package entity

type Farm struct {
	UID         string
	Name        string
	Description string
	Latitude    string
	Longitude   string
	Type        string
	CountryCode string
	CityCode    string
	IsActive    bool

	Reservoirs []Reservoir
}

// CreateFarm registers a new farm to Tania
func CreateFarm(name string, description string, latitude string, longitude string, farmType string, countryCode string, cityCode string) (Farm, error) {
	err := validateFarmName(name)
	if err != nil {
		return Farm{}, err
	}

	err = validateGeoLocation(latitude, longitude)
	if err != nil {
		return Farm{}, err
	}

	err = validateFarmType(farmType)
	if err != nil {
		return Farm{}, err
	}

	err = validateCountryCode(countryCode)
	if err != nil {
		return Farm{}, err
	}

	err = validateCityCode(countryCode, cityCode)
	if err != nil {
		return Farm{}, err
	}

	return Farm{
		Name:        name,
		Description: description,
		Latitude:    latitude,
		Longitude:   longitude,
		Type:        farmType,
		CountryCode: countryCode,
		CityCode:    cityCode,
		IsActive:    true,
	}, nil
}

func (f *Farm) AddReservoir(res Reservoir) error {
	found := f.IsReservoirAdded(res.Name)

	if found {
		return FarmError{FarmErrorReservoirAlreadyAdded}
	}

	f.Reservoirs = append(f.Reservoirs, res)

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

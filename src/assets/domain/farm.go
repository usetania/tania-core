// Package domain provides the operation that farm holder can do
// to their farm
package domain

import (
	"time"

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
	CreatedDate time.Time `json:"created_date"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

func (state *Farm) TrackChange(event interface{}) {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	state.Transition(event)
}

func (state *Farm) Transition(event interface{}) {
	switch e := event.(type) {
	case FarmCreated:
		state.UID = e.UID
		state.Name = e.Name
		state.Type = e.Type
		state.Latitude = e.Latitude
		state.Longitude = e.Longitude
		state.CountryCode = e.CountryCode
		state.CityCode = e.CityCode
		state.IsActive = e.IsActive
		state.CreatedDate = e.CreatedDate

	case FarmGeolocationChanged:
		state.Latitude = e.Latitude
		state.Longitude = e.Longitude

	case FarmRegionChanged:
		state.CountryCode = e.CountryCode
		state.CityCode = e.CityCode

	}
}

// CreateFarm registers a new farm to Tania
func CreateFarm(name, farmType, latitude, longitude, countryCode, cityCode string) (*Farm, error) {
	err := validateFarmName(name)
	if err != nil {
		return nil, err
	}

	err = validateFarmType(farmType)
	if err != nil {
		return nil, err
	}

	err = validateGeoLocation(latitude, longitude)
	if err != nil {
		return nil, err
	}

	err = validateCountryCode(countryCode)
	if err != nil {
		return nil, err
	}

	err = validateCityCode(countryCode, cityCode)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	initial := &Farm{}

	initial.TrackChange(FarmCreated{
		UID:         uid,
		Name:        name,
		Type:        farmType,
		Latitude:    latitude,
		Longitude:   longitude,
		CountryCode: countryCode,
		CityCode:    cityCode,
		IsActive:    true,
		CreatedDate: time.Now(),
	})

	return initial, nil
}

// ChangeGeoLocation changes the geolocation of a farm
func (f *Farm) ChangeGeoLocation(latitude, longitude string) error {
	err := validateGeoLocation(latitude, longitude)
	if err != nil {
		return err
	}

	f.TrackChange(FarmGeolocationChanged{
		FarmUID:   f.UID,
		Latitude:  latitude,
		Longitude: longitude,
	})

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

	f.TrackChange(FarmRegionChanged{
		FarmUID:     f.UID,
		CountryCode: countryCode,
		CityCode:    cityCode,
	})

	return nil
}

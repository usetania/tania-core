// Package domain provides the operation that farm holder can do
// to their farm
package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Farm struct {
	UID         uuid.UUID `json:"uid"`
	Name        string    `json:"name"`
	Latitude    string    `json:"latitude"`
	Longitude   string    `json:"longitude"`
	Type        string    `json:"type"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`

	// Events
	Version            int
	UncommittedChanges []interface{}
}

type FarmService interface {
	GetCountryNameByCode() string
}

type FarmCountry struct {
	Code string `json:"country_code"`
	Name string `json:"country_name"`
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
		state.Country = e.Country
		state.City = e.City
		state.IsActive = e.IsActive
		state.CreatedDate = e.CreatedDate

	case FarmNameChanged:
		state.Name = e.Name

	case FarmTypeChanged:
		state.Type = e.Type

	case FarmGeolocationChanged:
		state.Latitude = e.Latitude
		state.Longitude = e.Longitude

	case FarmRegionChanged:
		state.Country = e.Country
		state.City = e.City

	}
}

// CreateFarm registers a new farm to Tania
func CreateFarm(name, farmType, latitude, longitude, country, city string) (*Farm, error) {
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

	err = validateCountry(country)
	if err != nil {
		return nil, err
	}

	err = validateCity(city)
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
		Country:     country,
		City:        city,
		IsActive:    true,
		CreatedDate: time.Now(),
	})

	return initial, nil
}

func (f *Farm) ChangeName(name string) error {
	err := validateFarmName(name)
	if err != nil {
		return err
	}

	f.TrackChange(FarmNameChanged{
		FarmUID: f.UID,
		Name:    name,
	})

	return nil
}

func (f *Farm) ChangeType(farmType string) error {
	err := validateFarmType(farmType)
	if err != nil {
		return err
	}

	f.TrackChange(FarmTypeChanged{
		FarmUID: f.UID,
		Type:    farmType,
	})

	return nil
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
func (f *Farm) ChangeRegion(country, city string) error {
	err := validateCountry(country)
	if err != nil {
		return err
	}

	err = validateCity(city)
	if err != nil {
		return err
	}

	f.TrackChange(FarmRegionChanged{
		FarmUID: f.UID,
		Country: country,
		City:    city,
	})

	return nil
}

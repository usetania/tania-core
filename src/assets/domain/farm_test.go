package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFarm(t *testing.T) {
	// Given

	// When
	farm, err := CreateFarm("My Farm 1", FarmTypeOrganic, "10.00", "11.00", "Indonesia", "Jakarta")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "My Farm 1", farm.Name)

	event, ok := farm.UncommittedChanges[0].(FarmCreated)
	assert.True(t, ok)
	assert.Equal(t, farm.UID, event.UID)
}

func TestInvalidCreateFarm(t *testing.T) {
	var tests = []struct {
		name               string
		latitude           string
		longitude          string
		farmType           string
		countryCode        string
		cityCode           string
		expectedCreateFarm error
	}{
		{"My Farm Family", "-90.000", "-180.000", "organic", "Indonesia", "Jakarta", nil},
		{"", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorNameEmptyCode}},
		{"My Farm Family", "-90.000", "-180.000", "wrongtype", "Indonesia", "Jakarta", FarmError{FarmErrorInvalidFarmTypeCode}},
		{"My Farm Family", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorInvalidCountry}},
		{"My Farm Family", "-90.000", "-180.000", "organic", "Indonesia", "", FarmError{FarmErrorInvalidCity}},
	}

	for _, test := range tests {
		_, err := CreateFarm(
			test.name,
			test.farmType,
			test.latitude,
			test.longitude,
			test.countryCode,
			test.cityCode,
		)

		assert.Equal(t, test.expectedCreateFarm, err)
	}
}

func TestChangeGeolocation(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("my farm", "organic", "90.000", "100.000", "ID", "JK")
	latitude := "10.00"
	longitude := "11.00"

	// When
	geoErr := farm.ChangeGeoLocation(latitude, longitude)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, geoErr)

	assert.Equal(t, latitude, farm.Latitude)
	assert.Equal(t, longitude, farm.Longitude)

	event, ok := farm.UncommittedChanges[1].(FarmGeolocationChanged)
	assert.True(t, ok)
	assert.Equal(t, farm.UID, event.FarmUID)
	assert.Equal(t, farm.Latitude, event.Latitude)
}

func TestChangeRegion(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("my farm", "organic", "90.000", "100.000", "Indonesia", "Jakarta")
	country := "Australia"
	city := "Sydney"

	// When
	regErr := farm.ChangeRegion(country, city)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, regErr)

	assert.Equal(t, country, farm.Country)
	assert.Equal(t, city, farm.City)

	event, ok := farm.UncommittedChanges[1].(FarmRegionChanged)
	assert.True(t, ok)
	assert.Equal(t, farm.UID, event.FarmUID)
	assert.Equal(t, farm.Country, event.Country)
}

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFarm(t *testing.T) {
	var tests = []struct {
		name               string
		latitude           string
		longitude          string
		farmType           string
		countryCode        string
		cityCode           string
		expectedCreateFarm error
	}{
		{"My Farm Family", "-90.000", "-180.000", "organic", "ID", "JK", nil},
		{"", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorNameEmptyCode}},
		{"My Farm Family", "-90.000", "-180.000", "wrongtype", "ID", "JK", FarmError{FarmErrorInvalidFarmTypeCode}},
		{"My Farm Family", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorInvalidCountryCode}},
		{"My Farm Family", "-90.000", "-180.000", "organic", "ID", "Jakarta", FarmError{FarmErrorInvalidCityCode}},
	}

	for _, test := range tests {
		farm, err := CreateFarm(
			test.name,
			test.farmType,
			test.latitude,
			test.longitude,
			test.countryCode,
			test.cityCode,
		)

		assert.Equal(t, test.expectedCreateFarm, err)
		assert.NotEqual(t, Farm{}, farm)
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
}

func TestChangeRegion(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("my farm", "organic", "90.000", "100.000", "ID", "JK")
	countryCode := "AU"
	cityCode := "QLD"

	// When
	regErr := farm.ChangeRegion(countryCode, cityCode)

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, regErr)

	assert.Equal(t, countryCode, farm.CountryCode)
	assert.Equal(t, cityCode, farm.CityCode)
}

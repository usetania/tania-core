package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFarm(t *testing.T) {
	// Given

	// When
	farm, err := CreateFarm("My Farm 1", FarmTypeOrganic, "10.00", "11.00", "ID", "JK")

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
		{"My Farm Family", "-90.000", "-180.000", "organic", "ID", "JK", nil},
		{"", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorNameEmptyCode}},
		{"My Farm Family", "-90.000", "-180.000", "wrongtype", "ID", "JK", FarmError{FarmErrorInvalidFarmTypeCode}},
		{"My Farm Family", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorInvalidCountryCode}},
		{"My Farm Family", "-90.000", "-180.000", "organic", "ID", "Jakarta", FarmError{FarmErrorInvalidCityCode}},
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

	event, ok := farm.UncommittedChanges[1].(FarmRegionChanged)
	assert.True(t, ok)
	assert.Equal(t, farm.UID, event.FarmUID)
	assert.Equal(t, farm.CountryCode, event.CountryCode)
}

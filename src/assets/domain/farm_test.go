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
		{"MyFarmFamily", "-90.000", "-180.000", "organic", "ID", "JK", nil},
		{"", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorNameEmptyCode}},
		{"MyFarmFamily", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorInvalidCountryCode}},
		{"MyFarmFamily", "-90.000", "-180.000", "organic", "ID", "Jakarta", FarmError{FarmErrorInvalidCityCode}},
	}

	for _, test := range tests {
		farm, err := CreateFarm(test.name, test.farmType, test.latitude, test.longitude, test.countryCode, test.cityCode)

		assert.Equal(t, test.expectedCreateFarm, err)
		assert.NotEqual(t, Farm{}, farm)
	}
}

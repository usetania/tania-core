package entity

import (
	"testing"
)

func TestCreateFarm(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name        string
		description string
		latitude    string
		longitude   string
		farmType    string
		countryCode string
		cityCode    string
		expected    error
	}{
		{"My Farm Family", "", "-90.000", "-180.000", "organic", "ID", "JK", nil},
		{"", "", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorEmptyNameCode}},
		{"My Farm Family", "", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorInvalidCountryCode}},
		{"My Farm Family", "", "-90.000", "-180.000", "organic", "ID", "Jakarta", FarmError{FarmErrorInvalidCityCode}},
	}

	for _, test := range tests {
		_, actual := CreateFarm(test.name, test.description, test.latitude, test.longitude, test.farmType, test.countryCode, test.cityCode)

		if actual != test.expected {
			t.Errorf("Expected be %v, got %v", test.expected, actual)
		}
	}
}

package entity

import (
	"testing"
)

func TestValidateFarmName(t *testing.T) {
	t.Parallel()
	// Given
	var tests = []struct {
		param    string
		expected error
	}{
		{"MyAwesomeFarm", nil},
		{"", FarmError{FarmErrorNameEmptyCode}},
		{"Mys", FarmError{FarmErrorNameNotEnoughCharacterCode}},
		{"My4m gre<>", FarmError{FarmErrorNameAlphanumericOnlyCode}},
	}

	for _, test := range tests {
		// When
		actual := validateFarmName(test.param)
		if actual != test.expected {
			t.Errorf("Expected (%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestValidGeoLocation(t *testing.T) {
	t.Parallel()

	// Given
	var tests = []struct {
		latitude  string
		longitude string
		expected  error
	}{
		{"", "", FarmError{FarmErrorInvalidLatitudeValueCode}},
		{"108.000", "-180.000", FarmError{FarmErrorInvalidLatitudeValueCode}},
		{"+99.9", "-180.000", FarmError{FarmErrorInvalidLatitudeValueCode}},

		{"-90.000", "", FarmError{FarmErrorInvalidLongitudeValueCode}},
		{"-90.000", "180.1", FarmError{FarmErrorInvalidLongitudeValueCode}},
		{"-90.000", "+382.3811", FarmError{FarmErrorInvalidLongitudeValueCode}},

		{"-90.000", "-180.000", nil},
		{"+90.000", "-180.000", nil},
		{"47.1231231", "-180.000", nil},

		{"-90.000", "+72.234", nil},
		{"-90", "23.1111111", nil},
	}

	for _, test := range tests {
		actual := validateGeoLocation(test.latitude, test.longitude)
		if actual != test.expected {
			t.Errorf("Expected latitude (%q) longitude(%q) to be %v, got %v", test.latitude, test.longitude, test.expected, actual)
		}
	}
}

func TestValidateFarmType(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected error
	}{
		{"organic", nil},
		{"or", FarmError{FarmErrorInvalidFarmTypeCode}},
		{"", FarmError{FarmErrorInvalidFarmTypeCode}},
	}

	for _, test := range tests {
		actual := validateFarmType(test.param)
		if actual != test.expected {
			t.Errorf("Expected (%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestValidateCountryCode(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected error
	}{
		{"ID", nil},
		{"J", FarmError{FarmErrorInvalidCountryCode}},
		{"", FarmError{FarmErrorInvalidCountryCode}},
	}

	for _, test := range tests {
		actual := validateCountryCode(test.param)
		if actual != test.expected {
			t.Errorf("Expected (%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestValidateCityCode(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		countryCode string
		cityCode    string
		expected    error
	}{
		{"ID", "JK", nil},
		{"ID", "HKG", FarmError{FarmErrorInvalidCityCode}},
		{"ID", "", FarmError{FarmErrorInvalidCityCode}},
	}

	for _, test := range tests {
		actual := validateCityCode(test.countryCode, test.cityCode)
		if actual != test.expected {
			t.Errorf("Expected (%q, %q) to be %v, got %v", test.countryCode, test.cityCode, test.expected, actual)
		}
	}
}

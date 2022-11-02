//nolint:testpackage
package domain

import (
	"errors"
	"testing"
)

func TestValidateFarmName(t *testing.T) {
	t.Parallel()
	// Given
	tests := []struct {
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
		if !errors.Is(actual, test.expected) {
			t.Errorf("Expected (%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestValidGeoLocation(t *testing.T) {
	t.Parallel()

	// Given
	tests := []struct {
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
		if !errors.Is(actual, test.expected) {
			t.Errorf(
				"Expected latitude (%q) longitude(%q) to be %v, got %v",
				test.latitude, test.longitude, test.expected, actual,
			)
		}
	}
}

func TestValidateFarmType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected error
	}{
		{"organic", nil},
		{"or", FarmError{FarmErrorInvalidFarmTypeCode}},
		{"", FarmError{FarmErrorInvalidFarmTypeCode}},
	}

	for _, test := range tests {
		actual := validateFarmType(test.param)
		if !errors.Is(actual, test.expected) {
			t.Errorf("Expected (%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestValidateRegion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		country  string
		expected error
	}{
		{"Indonesia", nil},
		{"", FarmError{FarmErrorInvalidCountry}},
	}

	for _, test := range tests {
		actual := validateCountry(test.country)
		if !errors.Is(actual, test.expected) {
			t.Errorf("Expected (%q) to be %v, got %v", test.country, test.expected, actual)
		}
	}
}

func TestValidateCityCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		city     string
		expected error
	}{
		{"Jakarta", nil},
		{"", FarmError{FarmErrorInvalidCity}},
	}

	for _, test := range tests {
		actual := validateCity(test.city)
		if !errors.Is(actual, test.expected) {
			t.Errorf("Expected (%q) to be %v, got %v", test.city, test.expected, actual)
		}
	}
}

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		{"MyFarmFamily", "", "-90.000", "-180.000", "organic", "ID", "JK", nil},
		{"", "", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorNameEmptyCode}},
		{"MyFarmFamily", "", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorInvalidCountryCode}},
		{"MyFarmFamily", "", "-90.000", "-180.000", "organic", "ID", "Jakarta", FarmError{FarmErrorInvalidCityCode}},
	}

	for _, test := range tests {
		_, actual := CreateFarm(test.name, test.description, test.latitude, test.longitude, test.farmType, test.countryCode, test.cityCode)

		if actual != test.expected {
			t.Errorf("Expected be %v, got %v", test.expected, actual)
		}
	}
}

func TestAddReservoirToFarm(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")
	reservoir1, _ := CreateReservoir(farm, "MyReservoir1")
	reservoir2, _ := CreateReservoir(farm, "MyReservoir2")

	// When
	err1 := farm.AddReservoir(reservoir1)

	// Then
	assert.Equal(t, nil, err1)
	assert.Equal(t, len(farm.Reservoirs), 1)

	// When
	err2 := farm.AddReservoir(reservoir2)

	// Then
	assert.Equal(t, nil, err2)
	assert.Equal(t, len(farm.Reservoirs), 2)
}

func TestInvalidAddReservoirToFarm(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")
	reservoir, _ := CreateReservoir(farm, "MyReservoir1")

	// When
	err1 := farm.AddReservoir(reservoir)
	err2 := farm.AddReservoir(reservoir)

	// Then
	assert.Equal(t, nil, err1)
	assert.Equal(t, FarmError{FarmErrorReservoirAlreadyAdded}, err2)
}

func TestIsReservoirAddedInFarm(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")
	reservoir, _ := CreateReservoir(farm, "MyReservoir1")
	farm.AddReservoir(reservoir)

	// When
	result1 := farm.IsReservoirAdded("MyReservoir1")

	// Then
	assert.Equal(t, true, result1)

	// When
	result2 := farm.IsHaveReservoir()

	// Then
	assert.Equal(t, true, result2)
}

func TestInvalidIsReservoirAddedInFarm(t *testing.T) {
	// Given
	farm, _ := CreateFarm("Farm 1", "This is our farm", "10.00", "11.00", FarmTypeOrganic, "ID", "ID")

	// When
	result1 := farm.IsReservoirAdded("MyReservoir")

	// Then
	assert.Equal(t, false, result1)

	// When
	result2 := farm.IsHaveReservoir()

	// Then
	assert.Equal(t, false, result2)

	// Given
	reservoir, _ := CreateReservoir(farm, "MyReservoir1")
	farm.AddReservoir(reservoir)

	// When
	result3 := farm.IsReservoirAdded("MyReservoir")

	// Then
	assert.Equal(t, false, result3)
}

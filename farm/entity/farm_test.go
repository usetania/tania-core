package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFarm(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name                 string
		latitude             string
		longitude            string
		farmType             string
		countryCode          string
		cityCode             string
		expectedCreateFarm   error
		expectedChangeGeo    error
		expectedChangeRegion error
	}{
		{"MyFarmFamily", "-90.000", "-180.000", "organic", "ID", "JK", nil, nil, nil},
		{"", "-90.000", "-180.000", "organic", "", "Jakarta", FarmError{FarmErrorNameEmptyCode}, nil, nil},
		{"MyFarmFamily", "-90.000", "-180.000", "organic", "", "Jakarta", nil, nil, FarmError{FarmErrorInvalidCountryCode}},
		{"MyFarmFamily", "-90.000", "-180.000", "organic", "ID", "Jakarta", nil, nil, FarmError{FarmErrorInvalidCityCode}},
	}

	for _, test := range tests {
		farm, err := CreateFarm(test.name, test.farmType)

		assert.Equal(t, test.expectedCreateFarm, err)

		// check farm error to avoid null pointer reference of farm
		if err == nil {
			err = farm.ChangeGeoLocation(test.latitude, test.longitude)

			assert.Equal(t, test.expectedChangeGeo, err)

			err = farm.ChangeRegion(test.countryCode, test.cityCode)

			assert.Equal(t, test.expectedChangeRegion, err)
		}
	}
}

func TestAddReservoirToFarm(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
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
	assert.Nil(t, farmErr)
	assert.Equal(t, nil, err2)
	assert.Equal(t, len(farm.Reservoirs), 2)
}

func TestInvalidAddReservoirToFarm(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
	reservoir, _ := CreateReservoir(farm, "MyReservoir1")

	// When
	err1 := farm.AddReservoir(reservoir)
	err2 := farm.AddReservoir(reservoir)

	// Then
	assert.Nil(t, farmErr)
	assert.Equal(t, nil, err1)
	assert.Equal(t, FarmError{FarmErrorReservoirAlreadyAdded}, err2)
}

func TestIsReservoirAddedInFarm(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
	reservoir, _ := CreateReservoir(farm, "MyReservoir1")
	farm.AddReservoir(reservoir)

	// When
	result1 := farm.IsReservoirAdded("MyReservoir1")

	// Then
	assert.Nil(t, farmErr)
	assert.Equal(t, true, result1)

	// When
	result2 := farm.IsHaveReservoir()

	// Then
	assert.Equal(t, true, result2)
}

func TestInvalidIsReservoirAddedInFarm(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)

	// When
	result1 := farm.IsReservoirAdded("MyReservoir")

	// Then
	assert.Nil(t, farmErr)
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

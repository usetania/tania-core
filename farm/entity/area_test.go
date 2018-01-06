package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArea(t *testing.T) {
	farm, err := CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	reservoir, err := CreateReservoir(farm, "MyRes1")
	if err != nil {
		assert.Nil(t, err)
	}

	var tests = []struct {
		Name                  string
		Size                  AreaUnit
		Type                  string
		Location              string
		Photo                 AreaPhoto
		Reservoir             Reservoir
		Farm                  Farm
		expectedAreaError     error
		exptectedSizeError    error
		expectedLocationError error
	}{
		{"MyArea1", SquareMeter{Value: 100}, "nursery", "indoor", AreaPhoto{}, reservoir, farm, nil, nil, nil},
		{"MyArea2", Hectare{Value: 5}, "growing", "outdoor", AreaPhoto{}, reservoir, farm, nil, nil, nil},
	}

	for _, test := range tests {
		area, err := CreateArea(test.Farm, test.Name, test.Type)

		assert.Equal(t, test.expectedAreaError, err)

		if err == nil {
			err = area.ChangeSize(test.Size)

			assert.Equal(t, test.exptectedSizeError, err)

			err = area.ChangeLocation(test.Location)

			assert.Equal(t, test.expectedLocationError, err)

			assert.NotNil(t, area.UID)
		}
	}
}

// func TestAddReservoirToFarm(t *testing.T) {
// 	// Given
// 	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
// 	reservoir1, _ := CreateReservoir(farm, "MyReservoir1")
// 	reservoir2, _ := CreateReservoir(farm, "MyReservoir2")

// 	// When
// 	err1 := farm.AddReservoir(&reservoir1)

// 	// Then
// 	assert.Equal(t, nil, err1)
// 	assert.Equal(t, len(farm.Reservoirs), 1)

// 	// When
// 	err2 := farm.AddReservoir(&reservoir2)

// 	// Then
// 	assert.Nil(t, farmErr)
// 	assert.Equal(t, nil, err2)
// 	assert.Equal(t, len(farm.Reservoirs), 2)
// }

// func TestInvalidAddReservoirToFarm(t *testing.T) {
// 	// Given
// 	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
// 	reservoir, _ := CreateReservoir(farm, "MyReservoir1")

// 	// When
// 	err1 := farm.AddReservoir(&reservoir)
// 	err2 := farm.AddReservoir(&reservoir)

// 	// Then
// 	assert.Nil(t, farmErr)
// 	assert.Equal(t, nil, err1)
// 	assert.Equal(t, FarmError{FarmErrorReservoirAlreadyAdded}, err2)
// }

// func TestIsReservoirAddedInFarm(t *testing.T) {
// 	// Given
// 	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)
// 	reservoir, _ := CreateReservoir(farm, "MyReservoir1")
// 	farm.AddReservoir(&reservoir)

// 	// When
// 	result1 := farm.IsReservoirAdded("MyReservoir1")

// 	// Then
// 	assert.Nil(t, farmErr)
// 	assert.Equal(t, true, result1)

// 	// When
// 	result2 := farm.IsHaveReservoir()

// 	// Then
// 	assert.Equal(t, true, result2)
// }

// func TestInvalidIsReservoirAddedInFarm(t *testing.T) {
// 	// Given
// 	farm, farmErr := CreateFarm("Farm1", FarmTypeOrganic)

// 	// When
// 	result1 := farm.IsReservoirAdded("MyReservoir")

// 	// Then
// 	assert.Nil(t, farmErr)
// 	assert.Equal(t, false, result1)

// 	// When
// 	result2 := farm.IsHaveReservoir()

// 	// Then
// 	assert.Equal(t, false, result2)

// 	// Given
// 	reservoir, _ := CreateReservoir(farm, "MyReservoir1")
// 	farm.AddReservoir(&reservoir)

// 	// When
// 	result3 := farm.IsReservoirAdded("MyReservoir")

// 	// Then
// 	assert.Equal(t, false, result3)
// }

package domain

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

func TestCreateRemoveNote(t *testing.T) {
	// Given
	farm, farmErr := CreateFarm("MyFarm1", "organic")
	area, areaErr := CreateArea(farm, "Area1", "nursery")

	// When
	area.AddNewNote("This is my new note")

	// Then
	assert.Nil(t, farmErr)
	assert.Nil(t, areaErr)

	assert.Equal(t, 1, len(area.Notes))
	assert.Equal(t, "This is my new note", area.Notes[0].Content)
	assert.NotNil(t, area.Notes[0].CreatedDate)

	// When
	area.RemoveNote("This is my new note")

	assert.Equal(t, 0, len(area.Notes))
}

package domain

import (
	"testing"
)

func TestCreateArea(t *testing.T) {
	// farm, err := CreateFarm("MyFarm1", "organic", "10.000", "11.000", "ID", "JK")
	// if err != nil {
	// 	assert.Nil(t, err)
	// }

	// reservoir, err := CreateReservoir(farm, "MyRes1", )
	// if err != nil {
	// 	assert.Nil(t, err)
	// }

	// var tests = []struct {
	// 	Name                  string
	// 	Size                  AreaSize
	// 	Type                  string
	// 	Location              string
	// 	Photo                 AreaPhoto
	// 	Reservoir             Reservoir
	// 	Farm                  Farm
	// 	expectedAreaError     error
	// 	exptectedSizeError    error
	// 	expectedLocationError error
	// }{
	// 	{"MyArea1", AreaSize{Value: 100, Unit: AreaUnit{Symbol: SquareMeter}}, AreaTypeSeeding, AreaLocationIndoor, AreaPhoto{}, reservoir, farm, nil, nil, nil},
	// 	{"MyArea2", AreaSize{Value: 5, Unit: AreaUnit{Symbol: Hectare}}, AreaTypeGrowing, AreaLocationOutdoor, AreaPhoto{}, reservoir, farm, nil, nil, nil},
	// }

	// for _, test := range tests {
	// 	area, err := CreateArea(test.Farm, test.Name, test.Type)

	// 	assert.Equal(t, test.expectedAreaError, err)

	// 	if err == nil {
	// 		err = area.ChangeSize(test.Size)

	// 		assert.Equal(t, test.exptectedSizeError, err)

	// 		err = area.ChangeLocation(test.Location)

	// 		assert.Equal(t, test.expectedLocationError, err)

	// 		assert.NotNil(t, area.UID)
	// 	}
	// }
}

func TestAreaCreateRemoveNote(t *testing.T) {
	// // Given
	// farm, farmErr := CreateFarm("MyFarm1", "organic")
	// area, areaErr := CreateArea(farm, "Area1", AreaTypeSeeding)

	// // When
	// area.AddNewNote("This is my new note")

	// // Then
	// assert.Nil(t, farmErr)
	// assert.Nil(t, areaErr)

	// assert.Equal(t, 1, len(area.Notes))

	// uid := uuid.UUID{}
	// for k, v := range area.Notes {
	// 	assert.Equal(t, "This is my new note", v.Content)
	// 	assert.NotNil(t, v.CreatedDate)
	// 	uid = k
	// }

	// // When
	// area.RemoveNote(uid.String())

	// assert.Equal(t, 0, len(area.Notes))
}

package domain_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/usetania/tania-core/src/assets/domain"
)

type AreaServiceMock struct {
	mock.Mock
}

func (m *AreaServiceMock) FindFarmByID(uid uuid.UUID) (AreaFarmServiceResult, error) {
	args := m.Called(uid)

	return args.Get(0).(AreaFarmServiceResult), nil
}

func (m *AreaServiceMock) FindReservoirByID(uid uuid.UUID) (AreaReservoirServiceResult, error) {
	args := m.Called(uid)

	return args.Get(0).(AreaReservoirServiceResult), nil
}

func (m *AreaServiceMock) CountCropsByAreaID(areaUID uuid.UUID) (int, error) {
	args := m.Called(areaUID)

	return args.Get(0).(int), nil
}

type countCropsResult struct {
	AreaUID uuid.UUID
	Count   int
}

func TestCreateArea(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	farmResult := AreaFarmServiceResult{UID: farmUID}

	reservoirUID, _ := uuid.NewV4()
	reservoirResult := AreaReservoirServiceResult{UID: reservoirUID}

	countCropsResult := countCropsResult{
		AreaUID: uuid.UUID{},
		Count:   5,
	}
	areaService := mockAreaService(farmResult, reservoirResult, countCropsResult)

	// When
	area, err := CreateArea(
		areaService,
		farmUID,
		reservoirUID,
		"My Area 1",
		AreaTypeSeeding,
		AreaSize{Unit: GetAreaUnit(SquareMeter), Value: float32(10)},
		AreaLocationIndoor,
	)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "My Area 1", area.Name)

	event, ok := area.UncommittedChanges[0].(AreaCreated)
	assert.True(t, ok)
	assert.Equal(t, area.UID, event.UID)
}

func TestInvalidCreateArea(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	farmResult := AreaFarmServiceResult{UID: farmUID}

	reservoirUID, _ := uuid.NewV4()
	reservoirResult := AreaReservoirServiceResult{UID: reservoirUID}

	countCropsResult := countCropsResult{
		AreaUID: uuid.UUID{},
		Count:   5,
	}

	areaService := mockAreaService(farmResult, reservoirResult, countCropsResult)

	tests := []struct {
		Name          string
		Size          AreaSize
		Type          string
		Location      string
		ReservoirUID  uuid.UUID
		FarmUID       uuid.UUID
		ExpectedError error
	}{
		{
			Name:          "",
			Size:          AreaSize{Value: 100, Unit: GetAreaUnit(SquareMeter)},
			Type:          AreaTypeSeeding,
			Location:      AreaLocationIndoor,
			ReservoirUID:  reservoirUID,
			FarmUID:       farmUID,
			ExpectedError: AreaError{AreaErrorNameEmptyCode},
		},
		{
			Name:          "MyArea1",
			Size:          AreaSize{Value: 5, Unit: AreaUnit{Symbol: Hectare}},
			Type:          "WrongAreaType",
			Location:      AreaLocationOutdoor,
			ReservoirUID:  reservoirUID,
			FarmUID:       farmUID,
			ExpectedError: AreaError{AreaErrorInvalidAreaTypeCode},
		},
		{
			Name:          "MyArea1",
			Size:          AreaSize{Value: 5, Unit: AreaUnit{Symbol: Hectare}},
			Type:          AreaTypeSeeding,
			Location:      "WrongAreaLocation",
			ReservoirUID:  reservoirUID,
			FarmUID:       farmUID,
			ExpectedError: AreaError{Code: AreaErrorInvalidAreaLocationCode},
		},
	}

	for _, test := range tests {
		_, err := CreateArea(areaService, test.FarmUID, test.ReservoirUID, test.Name, test.Type, test.Size, test.Location)

		assert.Equal(t, test.ExpectedError, err)
	}
}

func TestAreaCreateRemoveNote(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	farmResult := AreaFarmServiceResult{UID: farmUID}

	reservoirUID, _ := uuid.NewV4()
	reservoirResult := AreaReservoirServiceResult{UID: reservoirUID}

	countCropsResult := countCropsResult{
		AreaUID: uuid.UUID{},
		Count:   5,
	}

	areaService := mockAreaService(farmResult, reservoirResult, countCropsResult)

	area, areaErr := CreateArea(
		areaService,
		farmUID,
		reservoirUID,
		"My Area 1",
		AreaTypeSeeding,
		AreaSize{Unit: GetAreaUnit(SquareMeter), Value: float32(10)},
		AreaLocationIndoor,
	)

	noteContent := "This is my new note"

	// When
	noteErr := area.AddNewNote(noteContent)

	// Then
	assert.Nil(t, areaErr)
	assert.Nil(t, noteErr)

	assert.Equal(t, 1, len(area.Notes))

	uid := uuid.UUID{}

	for k, v := range area.Notes {
		assert.Equal(t, noteContent, v.Content)
		assert.NotNil(t, v.CreatedDate)

		uid = k
	}

	event1, ok := area.UncommittedChanges[1].(AreaNoteAdded)
	assert.True(t, ok)
	assert.Equal(t, area.UID, event1.AreaUID)

	// When
	noteErr = area.RemoveNote(uid)

	// Then
	assert.Nil(t, noteErr)
	assert.Equal(t, 0, len(area.Notes))

	event2, ok := area.UncommittedChanges[2].(AreaNoteRemoved)
	assert.True(t, ok)
	assert.Equal(t, area.UID, event2.AreaUID)
}

func TestAreaChangePhoto(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	farmResult := AreaFarmServiceResult{UID: farmUID}

	reservoirUID, _ := uuid.NewV4()
	reservoirResult := AreaReservoirServiceResult{UID: reservoirUID}

	countCropsResult := countCropsResult{
		AreaUID: uuid.UUID{},
		Count:   5,
	}

	areaService := mockAreaService(farmResult, reservoirResult, countCropsResult)

	area, areaErr := CreateArea(
		areaService,
		farmUID,
		reservoirUID,
		"My Area 1",
		AreaTypeSeeding,
		AreaSize{Unit: GetAreaUnit(SquareMeter), Value: float32(10)},
		AreaLocationIndoor,
	)

	photo := AreaPhoto{
		Filename: "myphoto.jpg",
		MimeType: "image/jpeg",
		Size:     1000,
		Width:    800,
		Height:   600,
	}

	// When
	photoErr := area.ChangePhoto(photo)

	// Then
	assert.Nil(t, areaErr)
	assert.Nil(t, photoErr)
	assert.Equal(t, area.Photo.Filename, photo.Filename)

	event, ok := area.UncommittedChanges[1].(AreaPhotoAdded)
	assert.True(t, ok)
	assert.Equal(t, area.UID, event.AreaUID)
	assert.Equal(t, photo.Filename, event.Filename)
}

func mockAreaService(results ...interface{}) *AreaServiceMock {
	areaServiceMock := new(AreaServiceMock)

	for _, v := range results {
		switch res := v.(type) {
		case AreaFarmServiceResult:
			areaServiceMock.On("FindFarmByID", res.UID).Return(res)
		case AreaReservoirServiceResult:
			areaServiceMock.On("FindReservoirByID", res.UID).Return(res)
		case countCropsResult:
			areaServiceMock.On("CountCropsByAreaID", res.AreaUID).Return(res.Count)
		}
	}

	return areaServiceMock
}

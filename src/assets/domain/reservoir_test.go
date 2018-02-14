package domain

import (
	"testing"

	"github.com/Tanibox/tania-server/src/helper/mathhelper"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ReservoirServiceMock struct {
	mock.Mock
}

func (m ReservoirServiceMock) FindFarmByID(uid uuid.UUID) (FarmServiceResult, error) {
	args := m.Called(uid)
	return args.Get(0).(FarmServiceResult), nil
}

func TestCreateReservoir(t *testing.T) {
	// Given
	reservoirServiceMock := new(ReservoirServiceMock)

	farmUID, _ := uuid.NewV4()
	farmServiceResult := FarmServiceResult{
		UID:  farmUID,
		Name: "My Farm 1",
	}
	reservoirServiceMock.On("FindFarmByID", farmUID).Return(farmServiceResult)

	// When
	reservoir, err := CreateReservoir(reservoirServiceMock, farmUID, "My Reservoir 1", BucketType, float32(10))

	// Then
	assert.Nil(t, err)
	assert.NotEqual(t, Reservoir{}, reservoir)
}

func TestInvalidCreateReservoir(t *testing.T) {
	// Given
	reservoirServiceMock := new(ReservoirServiceMock)

	farmUID, _ := uuid.NewV4()
	farmServiceResult := FarmServiceResult{
		UID:  farmUID,
		Name: "My Farm 1",
	}
	reservoirServiceMock.On("FindFarmByID", farmUID).Return(farmServiceResult)

	reservoirData := []struct {
		farmUID         uuid.UUID
		name            string
		waterSourceType string
		capacity        float32
		expected        ReservoirError
	}{
		{farmUID, "My<>Reserv", "BUCKET", float32(10), ReservoirError{ReservoirErrorNameAlphanumericOnlyCode}},
		{farmUID, "MyR", "TAP", float32(0), ReservoirError{ReservoirErrorNameNotEnoughCharacterCode}},
		{farmUID, "MyReservoirNameShouldNotBeMoreThanAHundredLongCharactersMyReservoirNameShouldNotBeMoreThanAHundredLongCharacters", "TAP", 0, ReservoirError{ReservoirErrorNameExceedMaximunCharacterCode}},
	}

	for _, data := range reservoirData {
		// When
		_, err := CreateReservoir(reservoirServiceMock, data.farmUID, data.name, data.waterSourceType, data.capacity)

		// Then
		assert.Equal(t, data.expected, err)
	}
}

func TestReservoirCreateRemoveNote(t *testing.T) {
	// Given
	reservoirServiceMock := new(ReservoirServiceMock)

	farmUID, _ := uuid.NewV4()
	farmServiceResult := FarmServiceResult{
		UID:  farmUID,
		Name: "My Farm 1",
	}
	reservoirServiceMock.On("FindFarmByID", farmUID).Return(farmServiceResult)

	reservoir, reservoirErr := CreateReservoir(reservoirServiceMock, farmUID, "MyReservoir", "BUCKET", float32(10))

	// When
	reservoir.AddNewNote("This is my new note")

	// Then
	assert.Nil(t, reservoirErr)

	assert.Equal(t, 1, len(reservoir.Notes))

	uid := uuid.UUID{}
	for k, v := range reservoir.Notes {
		assert.Equal(t, "This is my new note", v.Content)
		assert.NotNil(t, v.CreatedDate)
		uid = k
	}

	// When
	reservoir.RemoveNote(uid.String())

	// Then
	assert.Equal(t, 0, len(reservoir.Notes))
}

func TestCreateWaterSource(t *testing.T) {
	// Given

	// When
	bucket, err1 := CreateBucket(float32(100))
	tap, err2 := CreateTap()

	// Then
	assert.Nil(t, err1)
	assert.NotEqual(t, bucket, Bucket{})
	assert.InDelta(t, bucket.Capacity, float32(100), mathhelper.EPSILON)

	assert.Nil(t, err2)
	assert.Equal(t, tap, Tap{})
}

func TestInvalidCreateWaterSource(t *testing.T) {
	// When
	bucket1, err1 := CreateBucket(0)
	bucket2, err2 := CreateBucket(-1)

	// Then
	assert.Equal(t, ReservoirError{ReservoirErrorBucketCapacityInvalidCode}, err1)
	assert.Equal(t, Bucket{}, bucket1)

	assert.Equal(t, ReservoirError{ReservoirErrorBucketCapacityInvalidCode}, err2)
	assert.Equal(t, Bucket{}, bucket2)
}

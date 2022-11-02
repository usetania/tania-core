package domain_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/usetania/tania-core/src/assets/domain"
)

type ReservoirServiceMock struct {
	mock.Mock
}

func (m *ReservoirServiceMock) FindFarmByID(uid uuid.UUID) (ReservoirFarmServiceResult, error) {
	args := m.Called(uid)

	return args.Get(0).(ReservoirFarmServiceResult), nil
}

func mockReservoirService(farmUID uuid.UUID, farmName string) *ReservoirServiceMock {
	reservoirServiceMock := new(ReservoirServiceMock)

	reservoirFarmServiceResult := ReservoirFarmServiceResult{
		UID:  farmUID,
		Name: farmName,
	}
	reservoirServiceMock.On("FindFarmByID", farmUID).Return(reservoirFarmServiceResult)

	return reservoirServiceMock
}

func TestCreateReservoir(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	serviceMock := mockReservoirService(farmUID, "My Farm 1")

	// When
	reservoir, err := CreateReservoir(serviceMock, farmUID, "My Reservoir 1", BucketType, float32(10))

	// Then
	assert.Nil(t, err)
	assert.NotEqual(t, Reservoir{}, reservoir)

	event, ok := reservoir.UncommittedChanges[0].(ReservoirCreated)
	assert.True(t, ok)
	assert.Equal(t, reservoir.UID, event.UID)
}

func TestInvalidCreateReservoir(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	serviceMock := mockReservoirService(farmUID, "My Farm")

	reservoirData := []struct {
		farmUID         uuid.UUID
		name            string
		waterSourceType string
		capacity        float32
		expectedError   ReservoirError
	}{
		{farmUID, "My<>Reserv", BucketType, float32(10), ReservoirError{ReservoirErrorNameAlphanumericOnlyCode}},
		{farmUID, "MyR", TapType, float32(0), ReservoirError{ReservoirErrorNameNotEnoughCharacterCode}},
		{
			farmUID,
			"MyReservoirNameShouldNotBeMoreThanAHundredLongCharactersMyReservoirNameShouldNotBeMoreThanAHundredLongCharacters",
			TapType,
			0,
			ReservoirError{ReservoirErrorNameExceedMaximunCharacterCode},
		},
	}

	for _, data := range reservoirData {
		// When
		_, err := CreateReservoir(serviceMock, farmUID, data.name, data.waterSourceType, data.capacity)

		// Then
		assert.Equal(t, data.expectedError, err)
	}
}

func TestReservoirCreateRemoveNote(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	serviceMock := mockReservoirService(farmUID, "My Farm")

	noteContent := "This is my new note"

	reservoir, reservoirErr := CreateReservoir(serviceMock, farmUID, "MyReservoir", BucketType, float32(10))

	// When
	noteErr := reservoir.AddNewNote(noteContent)

	// Then
	assert.Nil(t, reservoirErr)
	assert.Nil(t, noteErr)

	assert.Equal(t, 1, len(reservoir.Notes))

	uid := uuid.UUID{}

	for k, v := range reservoir.Notes {
		assert.Equal(t, noteContent, v.Content)
		assert.NotEqual(t, uuid.UUID{}, v.UID)

		uid = k
	}

	event1, ok := reservoir.UncommittedChanges[1].(ReservoirNoteAdded)
	assert.True(t, ok)
	assert.Equal(t, reservoir.UID, event1.ReservoirUID)

	// When
	reservoir.RemoveNote(uid)

	// Then
	assert.Equal(t, 0, len(reservoir.Notes))

	event2, ok := reservoir.UncommittedChanges[2].(ReservoirNoteRemoved)
	assert.True(t, ok)
	assert.Equal(t, reservoir.UID, event2.ReservoirUID)
}

func TestCreateWaterSource(t *testing.T) {
	t.Parallel()
	// Given
	capacity := float32(100)

	// When
	bucket, err1 := CreateBucket(capacity)
	tap, err2 := CreateTap()

	// Then
	assert.Nil(t, err1)
	assert.NotEqual(t, bucket, Bucket{})
	assert.Equal(t, capacity, bucket.Capacity)
	assert.Equal(t, BucketType, bucket.Type())

	assert.Nil(t, err2)
	assert.Equal(t, tap, Tap{})
	assert.Equal(t, TapType, tap.Type())
}

func TestInvalidCreateWaterSource(t *testing.T) {
	t.Parallel()
	// When
	bucket1, err1 := CreateBucket(0)
	bucket2, err2 := CreateBucket(-1)

	// Then
	assert.Equal(t, ReservoirError{ReservoirErrorBucketCapacityInvalidCode}, err1)
	assert.Equal(t, Bucket{}, bucket1)

	assert.Equal(t, ReservoirError{ReservoirErrorBucketCapacityInvalidCode}, err2)
	assert.Equal(t, Bucket{}, bucket2)
}

func TestReservoirChangeWaterSource(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	serviceMock := mockReservoirService(farmUID, "My Farm")

	reservoirBucket, resBucketErr := CreateReservoir(serviceMock, farmUID, "MyReservoir Bucket", BucketType, float32(10))
	reservoirTap, resTapErr := CreateReservoir(serviceMock, farmUID, "MyReservoir Tap", TapType, 0)

	// When
	reservoirBucket.ChangeWaterSource(TapType, 0)
	reservoirTap.ChangeWaterSource(BucketType, float32(100))

	// Then
	assert.Nil(t, resBucketErr)
	assert.Nil(t, resTapErr)

	assert.Equal(t, TapType, reservoirBucket.WaterSource.Type())

	_, ok := reservoirBucket.WaterSource.(Tap)
	assert.True(t, ok)

	assert.Equal(t, BucketType, reservoirTap.WaterSource.Type())

	event1, ok := reservoirBucket.UncommittedChanges[1].(ReservoirWaterSourceChanged)
	assert.True(t, ok)
	assert.Equal(t, reservoirBucket.UID, event1.ReservoirUID)

	// Then
	bucket, ok := reservoirTap.WaterSource.(Bucket)
	assert.True(t, ok)
	assert.Equal(t, float32(100), bucket.Capacity)

	event2, ok := reservoirTap.UncommittedChanges[1].(ReservoirWaterSourceChanged)
	assert.True(t, ok)
	assert.Equal(t, reservoirTap.UID, event2.ReservoirUID)
}

func TestReservoirChangeName(t *testing.T) {
	t.Parallel()
	// Given
	farmUID, _ := uuid.NewV4()
	serviceMock := mockReservoirService(farmUID, "My Farm")

	res, resErr := CreateReservoir(serviceMock, farmUID, "My Reservoir", BucketType, float32(10))

	// When
	res.ChangeName("My Reservoir Changed")

	// Then
	assert.Nil(t, resErr)
	assert.Equal(t, "My Reservoir Changed", res.Name)

	event, ok := res.UncommittedChanges[1].(ReservoirNameChanged)
	assert.True(t, ok)
	assert.Equal(t, res.UID, event.ReservoirUID)
	assert.Equal(t, res.Name, event.Name)
}

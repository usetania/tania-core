package inmemory_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/repository/inmemory"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirServiceMock struct {
	mock.Mock
}

func (m *ReservoirServiceMock) FindFarmByID(uid uuid.UUID) (domain.ReservoirFarmServiceResult, error) {
	args := m.Called(uid)

	return args.Get(0).(domain.ReservoirFarmServiceResult), nil
}

func TestReservoirEventInMemorySave(t *testing.T) {
	t.Parallel()
	// Given
	done := make(chan bool)

	reservoirEventStorage := storage.CreateReservoirEventStorage()
	repo := inmemory.NewReservoirEventRepositoryInMemory(reservoirEventStorage)

	reservoirServiceMock := new(ReservoirServiceMock)

	farmUID, _ := uuid.NewV4()
	reservoirFarmServiceResult := domain.ReservoirFarmServiceResult{
		UID:  farmUID,
		Name: "My Farm 1",
	}
	reservoirServiceMock.On("FindFarmByID", farmUID).Return(reservoirFarmServiceResult)

	reservoir1, resErr1 := domain.CreateReservoir(reservoirServiceMock, farmUID, "MyReservoir1", "BUCKET", float32(10))
	reservoir2, resErr2 := domain.CreateReservoir(reservoirServiceMock, farmUID, "MyReservoir2", "TAP", float32(0))

	// When
	var err1, err2 error

	go func() {
		err1 = <-repo.Save(reservoir1.UID, reservoir1.Version, reservoir1.UncommittedChanges)
		err2 = <-repo.Save(reservoir2.UID, reservoir2.Version, reservoir2.UncommittedChanges)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, resErr1)
	assert.Nil(t, resErr2)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

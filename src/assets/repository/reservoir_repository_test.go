package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/src/assets/storage"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestReservoirInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	reservoirStorage := storage.ReservoirStorage{ReservoirMap: make(map[uuid.UUID]domain.Reservoir)}
	repo := NewReservoirRepositoryInMemory(&reservoirStorage)

	farm, farmErr := domain.CreateFarm("Farm1", domain.FarmTypeOrganic)

	reservoir1, _ := domain.CreateReservoir(farm, "MyReservoir1")

	reservoir2, _ := domain.CreateReservoir(farm, "MyReservoir2")
	bucket2, _ := domain.CreateBucket(100, 50)
	reservoir2.AttachBucket(bucket2)

	reservoir3, _ := domain.CreateReservoir(farm, "MyReservoir3")
	tap3, _ := domain.CreateTap()
	reservoir3.AttachTap(tap3)

	// When
	var err1, err2, err3 error
	go func() {
		err1 = <-repo.Save(&reservoir1)
		err2 = <-repo.Save(&reservoir2)
		err3 = <-repo.Save(&reservoir3)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
}

func TestReservoirInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)
	reservoirStorage := storage.ReservoirStorage{ReservoirMap: make(map[uuid.UUID]domain.Reservoir)}
	repo := NewReservoirRepositoryInMemory(&reservoirStorage)

	farm, farmErr1 := domain.CreateFarm("Farm1", domain.FarmTypeOrganic)
	reservoir1, _ := domain.CreateReservoir(farm, "MyReservoir1")
	reservoir2, _ := domain.CreateReservoir(farm, "MyReservoir2")
	reservoir3, _ := domain.CreateReservoir(farm, "MyReservoir3")

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&reservoir1)
		<-repo.Save(&reservoir2)
		<-repo.Save(&reservoir3)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.Reservoir)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)

	val1, ok := result.Result.([]domain.Reservoir)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 3)
	assert.Contains(t, val1[0].Name, "MyReservoir")

	val2, ok := foundOne.Result.(domain.Reservoir)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
	assert.Contains(t, val2.Name, "MyReservoir")
}

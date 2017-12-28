package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/stretchr/testify/assert"
)

func TestReservoirInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	repo := NewReservoirRepositoryInMemory()

	farm, farmErr := entity.CreateFarm("Farm1", entity.FarmTypeOrganic)

	reservoir1, _ := entity.CreateReservoir(farm, "MyReservoir1")

	reservoir2, _ := entity.CreateReservoir(farm, "MyReservoir2")
	bucket2, _ := entity.CreateBucket(100, 50)
	reservoir2.AttachBucket(&bucket2)

	reservoir3, _ := entity.CreateReservoir(farm, "MyReservoir3")
	tap3, _ := entity.CreateTap()
	reservoir3.AttachTap(&tap3)

	reservoir1.UID = GetRandomUID()
	reservoir2.UID = GetRandomUID()
	reservoir3.UID = GetRandomUID()

	// When
	var saveResult1, saveResult2, saveResult3 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&reservoir1)
		saveResult2 = <-repo.Save(&reservoir2)
		saveResult3 = <-repo.Save(&reservoir3)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr)

	assert.NotNil(t, saveResult1)
	assert.NotNil(t, saveResult2)
	assert.NotNil(t, saveResult3)
}

func TestReservoirInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	repo := NewReservoirRepositoryInMemory()

	farm, farmErr1 := entity.CreateFarm("Farm1", entity.FarmTypeOrganic)
	reservoir1, _ := entity.CreateReservoir(farm, "MyReservoir1")
	reservoir2, _ := entity.CreateReservoir(farm, "MyReservoir2")
	reservoir3, _ := entity.CreateReservoir(farm, "MyReservoir3")

	reservoir1.UID = GetRandomUID()
	reservoir2.UID = GetRandomUID()
	reservoir3.UID = GetRandomUID()

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&reservoir1)
		<-repo.Save(&reservoir2)
		<-repo.Save(&reservoir3)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]entity.Reservoir)
		foundOne = <-repo.FindByID(val[0].UID)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)

	val1, ok := result.Result.([]entity.Reservoir)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 3)
	assert.Contains(t, val1[0].Name, "MyReservoir")

	val2, ok := foundOne.Result.(entity.Reservoir)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
	assert.Contains(t, val2.Name, "MyReservoir")
}

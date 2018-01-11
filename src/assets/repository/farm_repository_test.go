package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/src/assets/storage"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestFarmInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	farmStorage := storage.FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm)}
	repo := NewFarmRepositoryInMemory(&farmStorage)

	farm1, farmErr1 := domain.CreateFarm("MyFarmFamily", "organic")
	farm2, farmErr2 := domain.CreateFarm("MySecondFarm", "organic")

	// When
	var err1, err2 error
	go func() {
		err1 = <-repo.Save(&farm1)
		err2 = <-repo.Save(&farm2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestFarmInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)
	farmStorage := storage.FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm)}
	repo := NewFarmRepositoryInMemory(&farmStorage)

	farm1, farmErr1 := domain.CreateFarm("Farm1", domain.FarmTypeOrganic)
	farm2, farmErr2 := domain.CreateFarm("Farm2", domain.FarmTypeOrganic)

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&farm1)
		<-repo.Save(&farm2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.Farm)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	val1, ok := result.Result.([]domain.Farm)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 2)
	assert.Contains(t, val1[0].Name, "Farm")

	val2, ok := foundOne.Result.(domain.Farm)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
	assert.Contains(t, val2.Name, "Farm")
}

func TestFarmInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)
	farmStorage := storage.FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm)}
	repo := NewFarmRepositoryInMemory(&farmStorage)

	farm1, farmErr1 := domain.CreateFarm("Farm1", domain.FarmTypeOrganic)
	farm2, farmErr2 := domain.CreateFarm("Farm2", domain.FarmTypeOrganic)

	var found1, found2 RepositoryResult
	go func() {
		// Given
		<-repo.Save(&farm1)
		<-repo.Save(&farm2)

		// When
		found1 = <-repo.FindByID(farm1.UID.String())
		found2 = <-repo.FindByID(farm2.UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	farmResult1 := found1.Result.(domain.Farm)
	assert.Equal(t, "Farm1", farmResult1.Name)

	farmResult2 := found2.Result.(domain.Farm)
	assert.Equal(t, "Farm2", farmResult2.Name)
}

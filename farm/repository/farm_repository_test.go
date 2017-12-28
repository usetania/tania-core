package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/stretchr/testify/assert"
)

func TestFarmInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	repo := NewFarmRepositoryInMemory()

	farm1, farmErr1 := entity.CreateFarm("MyFarmFamily", "organic")
	farm2, farmErr2 := entity.CreateFarm("MySecondFarm", "organic")

	farm1.UID = GetRandomUID()
	farm2.UID = GetRandomUID()

	// When
	var saveResult1, saveResult2 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&farm1)
		saveResult2 = <-repo.Save(&farm2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	assert.NotNil(t, saveResult1)
	assert.NotNil(t, saveResult2)
}

func TestFarmInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	repo := NewFarmRepositoryInMemory()

	farm1, farmErr1 := entity.CreateFarm("Farm1", entity.FarmTypeOrganic)
	farm2, farmErr2 := entity.CreateFarm("Farm2", entity.FarmTypeOrganic)

	farm1.UID = GetRandomUID()
	farm2.UID = GetRandomUID()

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&farm1)
		<-repo.Save(&farm2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]entity.Farm)
		foundOne = <-repo.FindByID(val[0].UID)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	val1, ok := result.Result.([]entity.Farm)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 2)
	assert.Contains(t, val1[0].Name, "Farm")

	val2, ok := foundOne.Result.(entity.Farm)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
	assert.Contains(t, val2.Name, "Farm")
}

func TestFarmInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)

	repo := NewFarmRepositoryInMemory()

	farm1, farmErr1 := entity.CreateFarm("Farm1", entity.FarmTypeOrganic)
	farm2, farmErr2 := entity.CreateFarm("Farm2", entity.FarmTypeOrganic)

	farm1.UID = GetRandomUID()
	farm2.UID = GetRandomUID()

	var result1, result2, found1, found2 RepositoryResult
	go func() {
		// Given
		result1 = <-repo.Save(&farm1)
		result2 = <-repo.Save(&farm2)

		// When
		uid1, _ := result1.Result.(string)
		found1 = <-repo.FindByID(uid1)

		uid2, _ := result2.Result.(string)
		found2 = <-repo.FindByID(uid2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	farmResult1 := found1.Result.(entity.Farm)
	assert.Equal(t, "Farm1", farmResult1.Name)

	farmResult2 := found2.Result.(entity.Farm)
	assert.Equal(t, "Farm2", farmResult2.Name)
}

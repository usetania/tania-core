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

	farm1, _ := entity.CreateFarm("My Farm Family", "", "-90.000", "-180.000", "organic", "ID", "JK")
	farm2, _ := entity.CreateFarm("My Second Farm", "", "-90.000", "-180.000", "organic", "ID", "JK")

	// When
	var saveResult1, saveResult2, count1 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&farm1)
		saveResult2 = <-repo.Save(&farm2)

		count1 = <-repo.Count()
		done <- true
	}()

	// Then
	<-done
	assert.NotNil(t, saveResult1)

	assert.Equal(t, count1.Result, 2)
}

func TestFarmInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)

	repo := NewFarmRepositoryInMemory()

	farm1, _ := entity.CreateFarm("Farm1", "", "-90.000", "-180.000", entity.FarmTypeOrganic, "ID", "JK")
	farm2, _ := entity.CreateFarm("Farm2", "", "-90.000", "-180.000", entity.FarmTypeOrganic, "ID", "JK")

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
	farmResult1 := found1.Result.(entity.Farm)
	assert.Equal(t, "Farm1", farmResult1.Name)

	farmResult2 := found2.Result.(entity.Farm)
	assert.Equal(t, "Farm2", farmResult2.Name)
}

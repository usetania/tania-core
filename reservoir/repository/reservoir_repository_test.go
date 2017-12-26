package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/stretchr/testify/assert"
)

func TestInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	repo := NewReservoirRepositoryInMemory()

	reservoir1, _ := reservoir.CreateReservoir("My Reservoir 1")

	reservoir2, _ := reservoir.CreateReservoir("My Reservoir 2")
	bucket2, _ := reservoir.CreateBucket(100, 50)
	reservoir2.AttachBucket(&bucket2)

	reservoir3, _ := reservoir.CreateReservoir("My Reservoir 3")
	tap3, _ := reservoir.CreateTap()
	reservoir3.AttachTap(&tap3)

	// When
	var saveResult1, saveResult2, saveResult3, count1 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&reservoir1)
		saveResult2 = <-repo.Save(&reservoir2)
		saveResult3 = <-repo.Save(&reservoir3)

		count1 = <-repo.Count()
		done <- true
	}()

	// Then
	<-done
	assert.NotNil(t, saveResult1)
	assert.NotNil(t, saveResult2)
	assert.NotNil(t, saveResult3)
	assert.Equal(t, count1.Result, 3)
}

func TestInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	repo := NewReservoirRepositoryInMemory()

	reservoir1, _ := reservoir.CreateReservoir("My Reservoir 1")
	reservoir2, _ := reservoir.CreateReservoir("My Reservoir 2")
	reservoir3, _ := reservoir.CreateReservoir("My Reservoir 3")

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&reservoir1)
		<-repo.Save(&reservoir2)
		<-repo.Save(&reservoir3)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]reservoir.Reservoir)
		foundOne = <-repo.FindByID(val[0].UID)

		done <- true
	}()

	// Then
	<-done
	val1, ok := result.Result.([]reservoir.Reservoir)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 3)

	val2, ok := foundOne.Result.(reservoir.Reservoir)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
}

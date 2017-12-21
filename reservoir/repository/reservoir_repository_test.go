package repository

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/stretchr/testify/assert"
)

func TestInMemorySave(t *testing.T) {
	// Given
	var m sync.Mutex
	var repo ReservoirRepository = &ReservoirRepositoryInMemory{mutex: &m}

	reservoir1, _ := reservoir.CreateReservoir("My Reservoir 1", 8.1, 12.2, 30.0)

	reservoir2, _ := reservoir.CreateReservoir("My Reservoir 2", 8.1, 12.2, 30.0)
	bucket2, _ := reservoir.CreateBucket(100, 50)
	reservoir2.AttachBucket(&bucket2)

	reservoir3, _ := reservoir.CreateReservoir("My Reservoir 3", 8.1, 12.2, 30.0)
	tap3, _ := reservoir.CreateTap()
	reservoir3.AttachTap(&tap3)

	// When
	saveResult1 := repo.Save(&reservoir1)
	saveResult2 := repo.Save(&reservoir2)
	saveResult3 := repo.Save(&reservoir3)

	// Then
	assert.NotNil(t, saveResult1)
	assert.NotNil(t, saveResult2)
	assert.NotNil(t, saveResult3)
}

func TestInMemoryFindAll(t *testing.T) {
	// Given
	var m sync.Mutex
	var repo ReservoirRepository = &ReservoirRepositoryInMemory{mutex: &m}

	reservoir1, _ := reservoir.CreateReservoir("My Reservoir 1", 8.1, 12.2, 30.0)
	reservoir2, _ := reservoir.CreateReservoir("My Reservoir 2", 8.1, 12.2, 30.0)
	reservoir3, _ := reservoir.CreateReservoir("My Reservoir 3", 8.1, 12.2, 30.0)
	repo.Save(&reservoir1)
	repo.Save(&reservoir2)
	repo.Save(&reservoir3)

	// When
	results := repo.FindAll()

	// Then
	for val := range results {
		fmt.Println(val.Result)
	}

	assert.Equal(t, len(results), 3)
}

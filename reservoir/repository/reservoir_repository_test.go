package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/stretchr/testify/assert"
)

func TestInMemorySave(t *testing.T) {
	// Given
	var repo ReservoirRepository = &ReservoirRepositoryInMemory{}

	reservoir1, _ := reservoir.CreateReservoir("My Reservoir 1", 8.1, 12.2, 30.0)

	reservoir2, _ := reservoir.CreateReservoir("My Reservoir 2", 8.1, 12.2, 30.0)
	bucket2, _ := reservoir.CreateBucket(100, 50)
	reservoir2.AttachBucket(&bucket2)

	// When
	uid1, _ := repo.Save(reservoir1)
	uid2, _ := repo.Save(reservoir2)

	// Then
	assert.NotNil(t, uid1)

	assert.NotNil(t, uid2)
	assert.NotNil(t, uid2)
}

func TestInMemoryFindAll(t *testing.T) {
	// Given
	var repo ReservoirRepository = &ReservoirRepositoryInMemory{}

	reservoir1, _ := reservoir.CreateReservoir("My Reservoir 1", 8.1, 12.2, 30.0)
	reservoir2, _ := reservoir.CreateReservoir("My Reservoir 2", 8.1, 12.2, 30.0)
	reservoir3, _ := reservoir.CreateReservoir("My Reservoir 3", 8.1, 12.2, 30.0)
	repo.Save(reservoir1)
	repo.Save(reservoir2)
	repo.Save(reservoir3)

	// When
	reservoirs, _ := repo.FindAll()

	// Then
	assert.Equal(t, len(reservoirs), 3)
	assert.Equal(t, reservoirs[0].Name, reservoir1.Name)
}

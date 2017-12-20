package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/reservoir"
	"github.com/stretchr/testify/assert"
)

func TestInMemorySave(t *testing.T) {
	// Given
	repo := ReservoirRepositoryInMemory{}
	reservoir, _ := reservoir.CreateReservoir("My Reservoir 1", 8.1, 12.2, 30.0)

	// When
	uid := repo.Save(reservoir)

	// Then
	assert.NotNil(t, uid)
}

func TestInMemoryFindAll(t *testing.T) {
	// Given
	repo := ReservoirRepositoryInMemory{}
	reservoir, _ := reservoir.CreateReservoir("My Reservoir 1", 8.1, 12.2, 30.0)
	repo.Save(reservoir)

	// When
	reservoirs := repo.FindAll()

	// Then
	assert.Equal(t, len(reservoirs), 1)
}

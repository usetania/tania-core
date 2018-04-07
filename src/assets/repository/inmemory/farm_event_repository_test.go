package inmemory

import (
	"testing"

	"github.com/Tanibox/tania-core/src/assets/storage"

	"github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/stretchr/testify/assert"
)

func TestFarmEventInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)

	farmEventStorage := storage.CreateFarmEventStorage()
	repo := NewFarmEventRepositoryInMemory(farmEventStorage)

	farm1, farmErr1 := domain.CreateFarm("My Farm 1", "organic", "10.000", "11.000", "ID", "JK")
	farm2, farmErr2 := domain.CreateFarm("My Farm 2", "organic", "10.000", "11.000", "ID", "JK")

	// When
	var err1, err2 error
	go func() {
		err1 = <-repo.Save(farm1.UID, farm1.Version, farm1.UncommittedChanges)
		err2 = <-repo.Save(farm2.UID, farm2.Version, farm2.UncommittedChanges)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr1)
	assert.Nil(t, farmErr2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

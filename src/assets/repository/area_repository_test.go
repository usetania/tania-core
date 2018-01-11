package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/src/assets/storage"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAreaInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area)}
	repo := NewAreaRepositoryInMemory(&areaStorage)

	farm, err := domain.CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	area1, areaErr1 := domain.CreateArea(farm, "MyArea1", "nursery")
	area2, areaErr2 := domain.CreateArea(farm, "MyArea2", "growing")

	// When
	var err1, err2 error
	go func() {
		err1 = <-repo.Save(&area1)
		err2 = <-repo.Save(&area2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, areaErr1)
	assert.Nil(t, areaErr2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestAreaInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area)}
	repo := NewAreaRepositoryInMemory(&areaStorage)

	farm, err := domain.CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	area1, areaErr1 := domain.CreateArea(farm, "MyArea1", "nursery")
	area2, areaErr2 := domain.CreateArea(farm, "MyArea2", "growing")

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&area1)
		<-repo.Save(&area2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.Area)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, areaErr1)
	assert.Nil(t, areaErr2)

	val1, ok := result.Result.([]domain.Area)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 2)
	assert.Contains(t, val1[0].Name, "Area")

	val2, ok := foundOne.Result.(domain.Area)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
	assert.Contains(t, val2.Name, "Area")
}

func TestAreaInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)

	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area)}
	repo := NewAreaRepositoryInMemory(&areaStorage)

	farm, err := domain.CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	area1, areaErr1 := domain.CreateArea(farm, "MyArea1", "nursery")
	area2, areaErr2 := domain.CreateArea(farm, "MyArea2", "growing")

	var found1, found2 RepositoryResult
	go func() {
		// Given
		<-repo.Save(&area1)
		<-repo.Save(&area2)

		// When
		found1 = <-repo.FindByID(area1.UID.String())
		found2 = <-repo.FindByID(area2.UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, areaErr1)
	assert.Nil(t, areaErr2)

	areaResult1 := found1.Result.(domain.Area)
	assert.Equal(t, "MyArea1", areaResult1.Name)

	areaResult2 := found2.Result.(domain.Area)
	assert.Equal(t, "MyArea2", areaResult2.Name)
}

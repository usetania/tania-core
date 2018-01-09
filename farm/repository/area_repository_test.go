package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/farm/storage"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAreaInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]entity.Area)}
	repo := NewAreaStorageInMemory(&areaStorage)

	farm, err := entity.CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	area1, areaErr1 := entity.CreateArea(farm, "MyArea1", "nursery")
	area2, areaErr2 := entity.CreateArea(farm, "MyArea2", "growing")

	// When
	var saveResult1, saveResult2 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&area1)
		saveResult2 = <-repo.Save(&area2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, areaErr1)
	assert.Nil(t, areaErr2)

	assert.NotNil(t, saveResult1)
	assert.NotNil(t, saveResult2)
}

func TestAreaInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]entity.Area)}
	repo := NewAreaStorageInMemory(&areaStorage)

	farm, err := entity.CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	area1, areaErr1 := entity.CreateArea(farm, "MyArea1", "nursery")
	area2, areaErr2 := entity.CreateArea(farm, "MyArea2", "growing")

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&area1)
		<-repo.Save(&area2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]entity.Area)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, areaErr1)
	assert.Nil(t, areaErr2)

	val1, ok := result.Result.([]entity.Area)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(val1), 2)
	assert.Contains(t, val1[0].Name, "Area")

	val2, ok := foundOne.Result.(entity.Area)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
	assert.Contains(t, val2.Name, "Area")
}

func TestAreaInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)

	areaStorage := storage.AreaStorage{AreaMap: make(map[uuid.UUID]entity.Area)}
	repo := NewAreaStorageInMemory(&areaStorage)

	farm, err := entity.CreateFarm("MyFarm1", "organic")
	if err != nil {
		assert.Nil(t, err)
	}

	area1, areaErr1 := entity.CreateArea(farm, "MyArea1", "nursery")
	area2, areaErr2 := entity.CreateArea(farm, "MyArea2", "growing")

	var result1, result2, found1, found2 RepositoryResult
	go func() {
		// Given
		result1 = <-repo.Save(&area1)
		result2 = <-repo.Save(&area2)

		// When
		uid1, _ := result1.Result.(uuid.UUID)
		found1 = <-repo.FindByID(uid1.String())

		uid2, _ := result2.Result.(uuid.UUID)
		found2 = <-repo.FindByID(uid2.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, areaErr1)
	assert.Nil(t, areaErr2)

	areaResult1 := found1.Result.(entity.Area)
	assert.Equal(t, "MyArea1", areaResult1.Name)

	areaResult2 := found2.Result.(entity.Area)
	assert.Equal(t, "MyArea2", areaResult2.Name)
}

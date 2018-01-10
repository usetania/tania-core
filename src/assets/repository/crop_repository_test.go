package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/src/assets/storage"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCropInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop)}
	repo := NewCropRepositoryInMemory(&cropStorage)

	farm, farmErr := domain.CreateFarm("MyFarm", "organic")
	area, areaErr := domain.CreateArea(farm, "MyArea", "nursery")
	crop1, cropErr1 := domain.CreateCropBatch(area)
	crop2, cropErr2 := domain.CreateCropBatch(area)

	// When
	var saveResult1, saveResult2 RepositoryResult
	go func() {
		<-repo.Save(&crop1)
		<-repo.Save(&crop2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr)
	assert.Nil(t, areaErr)
	assert.Nil(t, cropErr1)
	assert.Nil(t, cropErr2)

	assert.Nil(t, saveResult1.Error)
	assert.Nil(t, saveResult2.Error)
}

func TestCropInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)
	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop)}
	repo := NewCropRepositoryInMemory(&cropStorage)

	farm, farmErr := domain.CreateFarm("MyFarm", "organic")
	area, areaErr := domain.CreateArea(farm, "MyArea", "nursery")
	crop1, cropErr1 := domain.CreateCropBatch(area)
	crop2, cropErr2 := domain.CreateCropBatch(area)

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&crop1)
		<-repo.Save(&crop2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.Crop)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr)
	assert.Nil(t, areaErr)
	assert.Nil(t, cropErr1)
	assert.Nil(t, cropErr2)

	val1, ok := result.Result.([]domain.Crop)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val1))

	val2, ok := foundOne.Result.(domain.Crop)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
}

func TestCropInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)
	cropStorage := storage.CropStorage{CropMap: make(map[uuid.UUID]domain.Crop)}
	repo := NewCropRepositoryInMemory(&cropStorage)

	farm, farmErr := domain.CreateFarm("MyFarm", "organic")
	area, areaErr := domain.CreateArea(farm, "MyArea", "nursery")
	crop1, cropErr1 := domain.CreateCropBatch(area)
	crop2, cropErr2 := domain.CreateCropBatch(area)

	var found1, found2 RepositoryResult
	go func() {
		// Given
		<-repo.Save(&crop1)
		<-repo.Save(&crop2)

		// When
		found1 = <-repo.FindByID(crop1.UID.String())
		found2 = <-repo.FindByID(crop2.UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, farmErr)
	assert.Nil(t, areaErr)
	assert.Nil(t, cropErr1)
	assert.Nil(t, cropErr2)

	cropResult1 := found1.Result.(domain.Crop)
	assert.NotNil(t, cropResult1.UID)
	assert.Equal(t, cropResult1.InitialArea, area)

	cropResult2 := found2.Result.(domain.Crop)
	assert.NotNil(t, cropResult2.UID)
	assert.Equal(t, cropResult2.InitialArea, area)
}

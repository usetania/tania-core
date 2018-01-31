package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/src/assets/storage"
	deadlock "github.com/sasha-s/go-deadlock"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestInventoryMaterialInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)

	rwMutex := deadlock.RWMutex{}
	Storage := storage.MaterialStorage{MaterialMap: make(map[uuid.UUID]domain.Material), Lock: &rwMutex}
	repo := NewMaterialRepositoryInMemory(&Storage)

	mts1, errMat1 := domain.CreateMaterialTypeSeed(domain.PlantTypeVegetable)
	material1, errMat2 := domain.CreateMaterial("Bayam Lu Hsieh", "12", domain.MoneyEUR, mts1, 20, domain.MaterialUnitPackets)
	mts2, errMat3 := domain.CreateMaterialTypeSeed(domain.PlantTypeFruit)
	material2, errMat4 := domain.CreateMaterial("Orange Number One", "20", domain.MoneyEUR, mts2, 30, domain.MaterialUnitPackets)

	// When
	var err1, err2 error
	go func() {
		err1 = <-repo.Save(&material1)
		err2 = <-repo.Save(&material2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, errMat1)
	assert.Nil(t, errMat2)
	assert.Nil(t, errMat3)
	assert.Nil(t, errMat4)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
}

func TestMaterialInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)

	rwMutex := deadlock.RWMutex{}
	Storage := storage.MaterialStorage{MaterialMap: make(map[uuid.UUID]domain.Material), Lock: &rwMutex}
	repo := NewMaterialRepositoryInMemory(&Storage)

	mts1, errMat1 := domain.CreateMaterialTypeSeed(domain.PlantTypeVegetable)
	material1, errMat2 := domain.CreateMaterial("Bayam Lu Hsieh", "12", domain.MoneyEUR, mts1, 20, domain.MaterialUnitPackets)
	mts2, errMat3 := domain.CreateMaterialTypeSeed(domain.PlantTypeFruit)
	material2, errMat4 := domain.CreateMaterial("Orange Number One", "20", domain.MoneyEUR, mts2, 30, domain.MaterialUnitPackets)

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&material1)
		<-repo.Save(&material2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.Material)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, errMat1)
	assert.Nil(t, errMat2)
	assert.Nil(t, errMat3)
	assert.Nil(t, errMat4)

	val1, ok := result.Result.([]domain.Material)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val1))

	val2, ok := foundOne.Result.(domain.Material)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
}

func TestMaterialInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)

	rwMutex := deadlock.RWMutex{}
	Storage := storage.MaterialStorage{MaterialMap: make(map[uuid.UUID]domain.Material), Lock: &rwMutex}
	repo := NewMaterialRepositoryInMemory(&Storage)

	mts1, errMat1 := domain.CreateMaterialTypeSeed(domain.PlantTypeVegetable)
	material1, errMat2 := domain.CreateMaterial("Bayam Lu Hsieh", "12", domain.MoneyEUR, mts1, 20, domain.MaterialUnitPackets)
	mts2, errMat3 := domain.CreateMaterialTypeSeed(domain.PlantTypeFruit)
	material2, errMat4 := domain.CreateMaterial("Orange Number One", "20", domain.MoneyEUR, mts2, 30, domain.MaterialUnitPackets)

	var found1, found2 RepositoryResult
	go func() {
		// Given
		<-repo.Save(&material1)
		<-repo.Save(&material2)

		// When
		found1 = <-repo.FindByID(material1.UID.String())
		found2 = <-repo.FindByID(material2.UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, errMat1)
	assert.Nil(t, errMat2)
	assert.Nil(t, errMat3)
	assert.Nil(t, errMat4)

	mat1 := found1.Result.(domain.Material)
	assert.Equal(t, "Bayam Lu Hsieh", mat1.Name)

	mat2 := found2.Result.(domain.Material)
	assert.Equal(t, "Orange Number One", mat2.Name)
}

package repository

import (
	"testing"

	"github.com/Tanibox/tania-server/src/assets/storage"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestInventoryMaterialInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	inventoryStorage := storage.InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial)}
	repo := NewInventoryMaterialRepositoryInMemory(&inventoryStorage)

	inv1, invErr1 := domain.CreateInventoryMaterial(domain.Vegetable{}, "Sawi Putih")
	inv2, invErr2 := domain.CreateInventoryMaterial(domain.Fruit{}, "Apple")

	// When
	var saveResult1, saveResult2 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&inv1)
		saveResult2 = <-repo.Save(&inv2)

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, invErr1)
	assert.Nil(t, invErr2)

	assert.Nil(t, saveResult1.Error)
	assert.Nil(t, saveResult2.Error)
}

func TestInventoryMaterialInMemoryFindAll(t *testing.T) {
	// Given
	done := make(chan bool)
	inventoryStorage := storage.InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial)}
	repo := NewInventoryMaterialRepositoryInMemory(&inventoryStorage)

	inv1, invErr1 := domain.CreateInventoryMaterial(domain.Vegetable{}, "Sawi Putih")
	inv2, invErr2 := domain.CreateInventoryMaterial(domain.Fruit{}, "Apple")

	var result, foundOne RepositoryResult
	go func() {
		// Given
		<-repo.Save(&inv1)
		<-repo.Save(&inv2)

		// When
		result = <-repo.FindAll()

		val := result.Result.([]domain.InventoryMaterial)
		foundOne = <-repo.FindByID(val[0].UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, invErr1)
	assert.Nil(t, invErr2)

	val1, ok := result.Result.([]domain.InventoryMaterial)
	assert.Equal(t, ok, true)
	assert.Equal(t, 2, len(val1))

	val2, ok := foundOne.Result.(domain.InventoryMaterial)
	assert.Equal(t, ok, true)
	assert.Equal(t, val2.UID, val1[0].UID)
}

func TestInventoryMaterialInMemoryFindByID(t *testing.T) {
	// Given
	done := make(chan bool)
	inventoryStorage := storage.InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial)}
	repo := NewInventoryMaterialRepositoryInMemory(&inventoryStorage)

	inv1, invErr1 := domain.CreateInventoryMaterial(domain.Vegetable{}, "Sawi Putih")
	inv2, invErr2 := domain.CreateInventoryMaterial(domain.Fruit{}, "Apple")

	var found1, found2 RepositoryResult
	go func() {
		// Given
		<-repo.Save(&inv1)
		<-repo.Save(&inv2)

		// When
		found1 = <-repo.FindByID(inv1.UID.String())
		found2 = <-repo.FindByID(inv2.UID.String())

		done <- true
	}()

	// Then
	<-done
	assert.Nil(t, invErr1)
	assert.Nil(t, invErr2)

	invResult1 := found1.Result.(domain.InventoryMaterial)
	assert.Equal(t, domain.Vegetable{}, invResult1.PlantType)
	assert.Equal(t, "Sawi Putih", invResult1.Variety)

	invResult2 := found2.Result.(domain.InventoryMaterial)
	assert.Equal(t, domain.Fruit{}, invResult2.PlantType)
	assert.Equal(t, "Apple", invResult2.Variety)
}

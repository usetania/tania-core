package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type InventoryMaterialRepository interface {
	Save(val *domain.InventoryMaterial) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// InventoryMaterialRepositoryInMemory is in-memory InventoryRepository db implementation
type InventoryMaterialRepositoryInMemory struct {
	Storage *storage.InventoryMaterialStorage
}

func NewInventoryMaterialRepositoryInMemory(s *storage.InventoryMaterialStorage) InventoryMaterialRepository {
	return &InventoryMaterialRepositoryInMemory{Storage: s}
}

func (f *InventoryMaterialRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		inventories := []domain.InventoryMaterial{}
		for _, val := range f.Storage.InventoryMaterialMap {
			inventories = append(inventories, val)
		}

		result <- RepositoryResult{Result: inventories}

		close(result)
	}()

	return result
}

// Save is to save
func (f *InventoryMaterialRepositoryInMemory) Save(val *domain.InventoryMaterial) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.InventoryMaterialMap[val.UID] = *val

		result <- RepositoryResult{Error: nil}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (f *InventoryMaterialRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: f.Storage.InventoryMaterialMap[uid]}
	}()

	return result
}

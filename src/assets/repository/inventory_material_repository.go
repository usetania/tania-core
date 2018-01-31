package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type MaterialRepository interface {
	Save(val *domain.Material) <-chan error
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// MaterialRepositoryInMemory is in-memory Repository db implementation
type MaterialRepositoryInMemory struct {
	Storage *storage.MaterialStorage
}

func NewMaterialRepositoryInMemory(s *storage.MaterialStorage) MaterialRepository {
	return &MaterialRepositoryInMemory{Storage: s}
}

func (f *MaterialRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		inventories := []domain.Material{}
		for _, val := range f.Storage.MaterialMap {
			inventories = append(inventories, val)
		}

		result <- RepositoryResult{Result: inventories}

		close(result)
	}()

	return result
}

// Save is to save
func (f *MaterialRepositoryInMemory) Save(val *domain.Material) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.MaterialMap[val.UID] = *val

		result <- nil

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (f *MaterialRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: f.Storage.MaterialMap[uid]}
	}()

	return result
}

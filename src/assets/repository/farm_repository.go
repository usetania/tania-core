package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmRepository interface {
	Save(val *domain.Farm) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// FarmRepositoryInMemory is in-memory FarmRepository db implementation
type FarmRepositoryInMemory struct {
	Storage *storage.FarmStorage
}

func NewFarmRepositoryInMemory(s *storage.FarmStorage) FarmRepository {
	return &FarmRepositoryInMemory{Storage: s}
}

func (f *FarmRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		farms := []domain.Farm{}
		for _, val := range f.Storage.FarmMap {
			farms = append(farms, val)
		}

		result <- RepositoryResult{Result: farms}

		close(result)
	}()

	return result
}

// Save is to save
func (f *FarmRepositoryInMemory) Save(val *domain.Farm) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.FarmMap[val.UID] = *val

		result <- RepositoryResult{Error: nil}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (f *FarmRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: f.Storage.FarmMap[uid]}
	}()

	return result
}

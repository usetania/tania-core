package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type CropRepository interface {
	Save(val *domain.Crop) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// CropRepositoryInMemory is in-memory CropRepository db implementation
type CropRepositoryInMemory struct {
	Storage *storage.CropStorage
}

func NewCropRepositoryInMemory(s *storage.CropStorage) CropRepository {
	return &CropRepositoryInMemory{Storage: s}
}

func (f *CropRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		farms := []domain.Crop{}
		for _, val := range f.Storage.CropMap {
			farms = append(farms, val)
		}

		result <- RepositoryResult{Result: farms}

		close(result)
	}()

	return result
}

// Save is to save
func (f *CropRepositoryInMemory) Save(val *domain.Crop) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.CropMap[val.UID] = *val

		result <- RepositoryResult{Error: nil}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (f *CropRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: f.Storage.CropMap[uid]}
	}()

	return result
}

package repository

import (
	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropEventRepository interface {
	Save(uid uuid.UUID, event interface{}) <-chan error
	FindAll() <-chan RepositoryResult
	FindByID(uid uuid.UUID) <-chan RepositoryResult
}

type CropEventRepositoryInMemory struct {
	Storage *storage.CropEventStorage
}

func NewCropEventRepositoryInMemory(s *storage.CropEventStorage) CropEventRepository {
	return &CropEventRepositoryInMemory{Storage: s}
}

func NewCropBatchFromHistory(events []interface{}) *domain.Crop {
	state := &domain.Crop{}
	for _, event := range events {
		state.Transition(event)
		state.Version++
	}
	return state
}

func (f *CropEventRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		cropEvents := make(map[uuid.UUID][]interface{})
		for key, val := range f.Storage.CropEventMap {
			cropEvents[key] = append(cropEvents[key], val)
		}

		result <- RepositoryResult{Result: cropEvents}

		close(result)
	}()

	return result
}

// Save is to save
func (f *CropEventRepositoryInMemory) Save(uid uuid.UUID, event interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.CropEventMap[uid] = event

		result <- nil

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (f *CropEventRepositoryInMemory) FindByID(uid uuid.UUID) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		result <- RepositoryResult{Result: f.Storage.CropEventMap[uid]}
	}()

	return result
}

package inmemory

import (
	"github.com/Tanibox/tania-core/src/growth/repository"
	"github.com/Tanibox/tania-core/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropEventRepositoryInMemory struct {
	Storage *storage.CropEventStorage
}

func NewCropEventRepositoryInMemory(s *storage.CropEventStorage) repository.CropEventRepository {
	return &CropEventRepositoryInMemory{Storage: s}
}

// Save is to save
func (f *CropEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.CropEvents = append(f.Storage.CropEvents, storage.CropEvent{
				CropUID: uid,
				Version: latestVersion,
				Event:   v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

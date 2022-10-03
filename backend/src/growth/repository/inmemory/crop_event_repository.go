package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/repository"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropEventRepositoryInMemory struct {
	Storage *storage.CropEventStorage
}

func NewCropEventRepositoryInMemory(s *storage.CropEventStorage) repository.CropEvent {
	return &CropEventRepositoryInMemory{Storage: s}
}

// Save is to save.
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

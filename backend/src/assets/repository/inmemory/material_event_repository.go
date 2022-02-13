package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type MaterialEventRepositoryInMemory struct {
	Storage *storage.MaterialEventStorage
}

func NewMaterialEventRepositoryInMemory(s *storage.MaterialEventStorage) repository.MaterialEvent {
	return &MaterialEventRepositoryInMemory{Storage: s}
}

func (f *MaterialEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++

			f.Storage.MaterialEvents = append(f.Storage.MaterialEvents, storage.MaterialEvent{
				MaterialUID: uid,
				Version:     latestVersion,
				Event:       v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

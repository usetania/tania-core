package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaEventRepositoryInMemory struct {
	Storage *storage.AreaEventStorage
}

func NewAreaEventRepositoryInMemory(s *storage.AreaEventStorage) repository.AreaEvent {
	return &AreaEventRepositoryInMemory{Storage: s}
}

func (f *AreaEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++

			f.Storage.AreaEvents = append(f.Storage.AreaEvents, storage.AreaEvent{
				AreaUID: uid,
				Version: latestVersion,
				Event:   v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

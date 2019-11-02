package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/repository"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type MaterialEventRepositoryInMemory struct {
	Storage *storage.MaterialEventStorage
}

func NewMaterialEventRepositoryInMemory(s *storage.MaterialEventStorage) repository.MaterialEventRepository {
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

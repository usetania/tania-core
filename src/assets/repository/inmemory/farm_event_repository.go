package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/repository"
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/gofrs/uuid"
)

type FarmEventRepositoryInMemory struct {
	Storage *storage.FarmEventStorage
}

func NewFarmEventRepositoryInMemory(s *storage.FarmEventStorage) repository.FarmEventRepository {
	return &FarmEventRepositoryInMemory{Storage: s}
}

// Save is to save
func (f *FarmEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.FarmEvents = append(f.Storage.FarmEvents, storage.FarmEvent{
				FarmUID: uid,
				Version: latestVersion,
				Event:   v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

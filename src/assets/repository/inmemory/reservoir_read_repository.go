package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/repository"
	"github.com/Tanibox/tania-core/src/assets/storage"
)

type ReservoirReadRepositoryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirReadRepositoryInMemory(s *storage.ReservoirReadStorage) repository.ReservoirReadRepository {
	return &ReservoirReadRepositoryInMemory{Storage: s}
}

func (f *ReservoirReadRepositoryInMemory) Save(reservoirRead *storage.ReservoirRead) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.ReservoirReadMap[reservoirRead.UID] = *reservoirRead

		result <- nil

		close(result)
	}()

	return result
}

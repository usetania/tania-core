package inmemory

import (
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirReadRepositoryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirReadRepositoryInMemory(s *storage.ReservoirReadStorage) repository.ReservoirRead {
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

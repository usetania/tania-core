package repository

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type ReservoirReadRepository interface {
	Save(reservoirRead *storage.ReservoirRead) <-chan error
}

type ReservoirReadRepositoryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirReadRepositoryInMemory(s *storage.ReservoirReadStorage) ReservoirReadRepository {
	return &ReservoirReadRepositoryInMemory{Storage: s}
}

// Save is to save
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

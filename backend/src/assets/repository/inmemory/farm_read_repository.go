package inmemory

import (
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type FarmReadRepositoryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadRepositoryInMemory(s *storage.FarmReadStorage) repository.FarmRead {
	return &FarmReadRepositoryInMemory{Storage: s}
}

func (f *FarmReadRepositoryInMemory) Save(farmRead *storage.FarmRead) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.FarmReadMap[farmRead.UID] = *farmRead

		result <- nil

		close(result)
	}()

	return result
}

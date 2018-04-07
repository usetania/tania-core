package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/repository"
	"github.com/Tanibox/tania-core/src/assets/storage"
)

type FarmReadRepositoryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadRepositoryInMemory(s *storage.FarmReadStorage) repository.FarmReadRepository {
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

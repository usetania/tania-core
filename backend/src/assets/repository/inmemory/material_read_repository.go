package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/repository"
	"github.com/Tanibox/tania-core/src/assets/storage"
)

type MaterialReadRepositoryInMemory struct {
	Storage *storage.MaterialReadStorage
}

func NewMaterialReadRepositoryInMemory(s *storage.MaterialReadStorage) repository.MaterialReadRepository {
	return &MaterialReadRepositoryInMemory{Storage: s}
}

// Save is to save
func (f *MaterialReadRepositoryInMemory) Save(materialRead *storage.MaterialRead) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.MaterialReadMap[materialRead.UID] = *materialRead

		result <- nil

		close(result)
	}()

	return result
}

package repository

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type MaterialReadRepository interface {
	Save(materialRead *storage.MaterialRead) <-chan error
}

type MaterialReadRepositoryInMemory struct {
	Storage *storage.MaterialReadStorage
}

func NewMaterialReadRepositoryInMemory(s *storage.MaterialReadStorage) MaterialReadRepository {
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

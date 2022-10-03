package inmemory

import (
	"github.com/usetania/tania-core/src/growth/repository"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropReadRepositoryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropReadRepositoryInMemory(s *storage.CropReadStorage) repository.CropRead {
	return &CropReadRepositoryInMemory{Storage: s}
}

// Save is to save.
func (f *CropReadRepositoryInMemory) Save(cropRead *storage.CropRead) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.CropReadMap[cropRead.UID] = *cropRead

		result <- nil

		close(result)
	}()

	return result
}

package repository

import (
	"github.com/Tanibox/tania-server/src/growth/storage"
)

type CropListRepository interface {
	Save(cropList *storage.CropList) <-chan error
}

type CropListRepositoryInMemory struct {
	Storage *storage.CropListStorage
}

func NewCropListRepositoryInMemory(s *storage.CropListStorage) CropListRepository {
	return &CropListRepositoryInMemory{Storage: s}
}

// Save is to save
func (f *CropListRepositoryInMemory) Save(cropList *storage.CropList) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.CropListMap[cropList.UID] = *cropList

		result <- nil

		close(result)
	}()

	return result
}

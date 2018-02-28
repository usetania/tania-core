package inmemory

import (
	"github.com/Tanibox/tania-server/src/growth/storage"
)

type CropActivityRepositoryInMemory struct {
	Storage *storage.CropActivityStorage
}

func NewCropActivityRepositoryInMemory(s *storage.CropActivityStorage) CropActivityRepository {
	return &CropActivityRepositoryInMemory{Storage: s}
}

// Save is to save
func (f *CropActivityRepositoryInMemory) Save(cropActivity *storage.CropActivity) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.CropActivityMap = append(f.Storage.CropActivityMap, *cropActivity)

		result <- nil

		close(result)
	}()

	return result
}

package inmemory

import "github.com/Tanibox/tania-server/src/growth/storage"

type CropReadRepositoryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropReadRepositoryInMemory(s *storage.CropReadStorage) CropReadRepository {
	return &CropReadRepositoryInMemory{Storage: s}
}

// Save is to save
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

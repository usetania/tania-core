package repository

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type AreaReadRepository interface {
	Save(areaRead *storage.AreaRead) <-chan error
}

type AreaReadRepositoryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaReadRepositoryInMemory(s *storage.AreaReadStorage) AreaReadRepository {
	return &AreaReadRepositoryInMemory{Storage: s}
}

// Save is to save
func (f *AreaReadRepositoryInMemory) Save(areaRead *storage.AreaRead) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		f.Storage.AreaReadMap[areaRead.UID] = *areaRead

		result <- nil

		close(result)
	}()

	return result
}

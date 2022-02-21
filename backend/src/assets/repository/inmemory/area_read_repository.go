package inmemory

import (
	"github.com/usetania/tania-core/src/assets/repository"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaReadRepositoryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaReadRepositoryInMemory(s *storage.AreaReadStorage) repository.AreaRead {
	return &AreaReadRepositoryInMemory{Storage: s}
}

// Save is to save.
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

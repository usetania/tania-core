package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type FarmReadRepository interface {
	Save(farmRead *storage.FarmRead) <-chan error
}

type FarmReadRepositoryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadRepositoryInMemory(s *storage.FarmReadStorage) FarmReadRepository {
	return &FarmReadRepositoryInMemory{Storage: s}
}

func NewFarmFromHistory(events []storage.FarmEvent) *domain.Farm {
	state := &domain.Farm{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}
	return state
}

// Save is to save
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

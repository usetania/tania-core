package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type MaterialEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type MaterialEventRepositoryInMemory struct {
	Storage *storage.MaterialEventStorage
}

func NewMaterialEventRepositoryInMemory(s *storage.MaterialEventStorage) MaterialEventRepository {
	return &MaterialEventRepositoryInMemory{Storage: s}
}

func NewMaterialFromHistory(events []storage.MaterialEvent) *domain.Material {
	state := &domain.Material{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}
	return state
}

// Save is to save
func (f *MaterialEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.MaterialEvents = append(f.Storage.MaterialEvents, storage.MaterialEvent{
				MaterialUID: uid,
				Version:     latestVersion,
				Event:       v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type AreaEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type AreaEventRepositoryInMemory struct {
	Storage *storage.AreaEventStorage
}

func NewAreaEventRepositoryInMemory(s *storage.AreaEventStorage) AreaEventRepository {
	return &AreaEventRepositoryInMemory{Storage: s}
}

func NewAreaFromHistory(events []storage.AreaEvent) *domain.Area {
	state := &domain.Area{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}
	return state
}

// Save is to save
func (f *AreaEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.AreaEvents = append(f.Storage.AreaEvents, storage.AreaEvent{
				AreaUID: uid,
				Version: latestVersion,
				Event:   v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

package repository

import (
	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
)

type TaskEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type TaskEventRepositoryInMemory struct {
	Storage *storage.TaskEventStorage
}

func NewTaskEventRepositoryInMemory(s *storage.TaskEventStorage) TaskEventRepository {
	return &TaskEventRepositoryInMemory{Storage: s}
}

func BuildTaskEventsFromHistory(events []storage.TaskEvent) *domain.Task {
	state := &domain.Task{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}
	return state
}

// Save is to save
func (f *TaskEventRepositoryInMemory) Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error {
	result := make(chan error)

	go func() {
		f.Storage.Lock.Lock()
		defer f.Storage.Lock.Unlock()

		for _, v := range events {
			latestVersion++
			f.Storage.TaskEvents = append(f.Storage.TaskEvents, storage.TaskEvent{
				TaskUID: uid,
				Version: latestVersion,
				Event:   v,
			})
		}

		result <- nil

		close(result)
	}()

	return result
}

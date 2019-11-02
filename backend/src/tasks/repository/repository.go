package repository

import (
	"github.com/Tanibox/tania-core/src/tasks/domain"
	"github.com/Tanibox/tania-core/src/tasks/storage"
	uuid "github.com/satori/go.uuid"
)

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

// EventWrapper is used to wrap the event interface with its struct name,
// so it will be easier to unmarshal later
type EventWrapper struct {
	EventName string
	EventData interface{}
}

type TaskEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

func BuildTaskFromEventHistory(taskService domain.TaskService, events []storage.TaskEvent) *domain.Task {
	state := &domain.Task{}
	for _, v := range events {
		state.Transition(taskService, v.Event)
		state.Version++
	}
	return state
}

type TaskReadRepository interface {
	Save(taskRead *storage.TaskRead) <-chan error
}

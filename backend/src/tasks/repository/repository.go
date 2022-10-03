package repository

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/storage"
)

// Result is a struct to wrap repository result
// so its easy to use it in channel.
type Result struct {
	Result interface{}
	Error  error
}

// EventWrapper is used to wrap the event interface with its struct name,
// so it will be easier to unmarshal later.
type EventWrapper struct {
	EventName string
	EventData interface{}
}

type TaskEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

func BuildTaskFromEventHistory(events []storage.TaskEvent) *domain.Task {
	state := &domain.Task{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

type TaskRead interface {
	Save(taskRead *storage.TaskRead) <-chan error
}

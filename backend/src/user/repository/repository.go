package repository

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/user/domain"
	"github.com/usetania/tania-core/src/user/storage"
)

// Result is a struct to wrap repository result
// so its easy to use it in channel.
type Result struct {
	Result interface{}
	Error  error
}

type UserEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type UserRead interface {
	Save(userRead *storage.UserRead) <-chan error
}

type UserAuth interface {
	Save(userAuth *storage.UserAuth) <-chan error
}

func NewUserFromHistory(events []storage.UserEvent) *domain.User {
	state := &domain.User{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

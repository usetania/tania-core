package repository

import (
	"github.com/Tanibox/tania-core/src/user/domain"
	"github.com/Tanibox/tania-core/src/user/storage"
	uuid "github.com/satori/go.uuid"
)

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

type UserEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type UserReadRepository interface {
	Save(userRead *storage.UserRead) <-chan error
}

type UserAuthRepository interface {
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

package query

import "github.com/gofrs/uuid"

type UserEvent interface {
	FindAllByID(userUID uuid.UUID) <-chan Result
}

type UserRead interface {
	FindByID(userUID uuid.UUID) <-chan Result
	FindByUsername(username string) <-chan Result
	FindByUsernameAndPassword(username, password string) <-chan Result
}

type UserAuth interface {
	FindByUserID(userUID uuid.UUID) <-chan Result
}

type Result struct {
	Result interface{}
	Error  error
}

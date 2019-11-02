package query

import uuid "github.com/satori/go.uuid"

type UserEventQuery interface {
	FindAllByID(userUID uuid.UUID) <-chan QueryResult
}

type UserReadQuery interface {
	FindByID(userUID uuid.UUID) <-chan QueryResult
	FindByUsername(username string) <-chan QueryResult
	FindByUsernameAndPassword(username, password string) <-chan QueryResult
}

type UserAuthQuery interface {
	FindByUserID(userUID uuid.UUID) <-chan QueryResult
}

type QueryResult struct {
	Result interface{}
	Error  error
}

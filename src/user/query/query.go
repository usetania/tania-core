package query

import uuid "github.com/satori/go.uuid"

type UserEventQuery interface {
	FindAllByID(farmUID uuid.UUID) <-chan QueryResult
}

type UserReadQuery interface {
	FindByID(userUID uuid.UUID) <-chan QueryResult
	FindByUsername(username string) <-chan QueryResult
}

type QueryResult struct {
	Result interface{}
	Error  error
}

package storage

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserEvent struct {
	UserUID     uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
}

type UserRead struct {
	UID         uuid.UUID `json:"uid"`
	Username    string    `json:"username"`
	Password    []byte    `json:"-"`
	CreatedDate time.Time `json:"created_date"`
	LastUpdated time.Time `json:"last_updated"`
}

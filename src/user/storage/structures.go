package storage

import (
	"time"

	"github.com/gofrs/uuid"
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

type UserAuth struct {
	UserUID      uuid.UUID `json:"uid"`
	AccessToken  string    `json:"access_token"`
	TokenExpires int       `json:"token_expires"`
	CreatedDate  time.Time `json:"created_date"`
	LastUpdated  time.Time `json:"last_updated"`
}

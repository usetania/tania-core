package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserCreated struct {
	UID         uuid.UUID
	Username    string
	Password    []byte
	CreatedDate time.Time
	LastUpdated time.Time
}

type PasswordChanged struct {
	UID         uuid.UUID
	NewPassword []byte
	DateChanged time.Time
}

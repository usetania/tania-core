package domain

import uuid "github.com/satori/go.uuid"

type UserCreated struct {
	UID      uuid.UUID
	Username string
	Password []byte
}

type PasswordChanged struct {
	UID         uuid.UUID
	NewPassword []byte
}

package domain

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID         uuid.UUID
	Username    string
	Password    []byte
	ClientID    string
	CreatedDate time.Time
	LastUpdated time.Time

	// Events
	Version            int
	UncommittedChanges []interface{}
}

type UserService interface {
	FindUserByUsername(username string) (UserServiceResult, error)
}

type UserServiceResult struct {
	UID      uuid.UUID
	Username string
}

func (u *User) TrackChange(event interface{}) {
	u.UncommittedChanges = append(u.UncommittedChanges, event)
	u.Transition(event)
}

func (u *User) Transition(event interface{}) {
	switch e := event.(type) {
	case UserCreated:
		u.UID = e.UID
		u.Username = e.Username
		u.Password = e.Password
		u.CreatedDate = e.CreatedDate
		u.LastUpdated = e.LastUpdated

	case PasswordChanged:
		u.Password = e.NewPassword
		u.LastUpdated = e.DateChanged
	}
}

func CreateUser(userService UserService, username, password, confirmPassword string) (*User, error) {
	if username == "" {
		return nil, UserError{UserErrorUsernameEmptyCode}
	}

	if len(username) < 5 {
		return nil, UserError{UserErrorInvalidUsernameLengthCode}
	}

	userResult, err := userService.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	if userResult.UID != (uuid.UUID{}) {
		return nil, UserError{UserErrorUsernameExistsCode}
	}

	err = validatePassword(password, confirmPassword)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	user := &User{
		UID:      uid,
		Username: username,
		Password: hash,
	}

	now := time.Now()

	user.TrackChange(UserCreated{
		UID:         uid,
		Username:    username,
		Password:    hash,
		CreatedDate: now,
		LastUpdated: now,
	})

	return user, nil
}

func (u *User) ChangePassword(oldPassword, newPassword, newConfirmPassword string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(oldPassword))
	if err != nil {
		return UserError{UserChangePasswordErrorWrongOldPasswordCode}
	}

	err = validatePassword(newPassword, newConfirmPassword)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %w", err)
	}

	u.TrackChange(PasswordChanged{
		UID:         u.UID,
		NewPassword: hash,
		DateChanged: time.Now(),
	})

	return nil
}

func (u *User) IsPasswordValid(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false, UserError{UserErrorWrongPasswordCode}
	}

	return true, nil
}

func validatePassword(password, confirmPassword string) error {
	if password == "" {
		return UserError{UserErrorPasswordEmptyCode}
	}

	if password != confirmPassword {
		return UserError{UserErrorPasswordConfirmationNotMatchCode}
	}

	return nil
}

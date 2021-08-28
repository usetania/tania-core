package domain

import (
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

func (state *User) TrackChange(event interface{}) {
	state.UncommittedChanges = append(state.UncommittedChanges, event)
	state.Transition(event)
}

func (state *User) Transition(event interface{}) {
	switch e := event.(type) {
	case UserCreated:
		state.UID = e.UID
		state.Username = e.Username
		state.Password = e.Password
		state.CreatedDate = e.CreatedDate
		state.LastUpdated = e.LastUpdated

	case PasswordChanged:
		state.Password = e.NewPassword
		state.LastUpdated = e.DateChanged

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
		return nil, err
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
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
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
		return err
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

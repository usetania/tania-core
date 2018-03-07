package domain

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID      uuid.UUID
	Username string
	Password []byte

	// Events
	Version            int
	UncommittedChanges []interface{}
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

	case PasswordChanged:
		state.Password = e.NewPassword

	}
}

func CreateUser(username, password, confirmPassword string) (*User, error) {
	if username == "" {
		return nil, UserError{UserErrorUsernameEmptyCode}
	}

	if len(username) < 5 {
		return nil, UserError{UserErrorInvalidUsernameLengthCode}
	}

	err := validatePassword(password, confirmPassword)
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

	user.TrackChange(UserCreated{
		UID:      uid,
		Username: username,
		Password: hash,
	})

	return user, nil
}

func (u *User) ChangePassword(oldPassword, newPassword, newConfirmPassword string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(oldPassword))
	if err != nil {
		return UserError{UserChangePasswordErrorWrongOldPassword}
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
		NewPassword: hash,
	})

	return nil
}

func (u *User) IsPasswordValid(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		return false, UserError{UserErrorWrongPassword}
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

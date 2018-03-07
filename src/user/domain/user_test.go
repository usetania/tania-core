package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Given

	// When
	user, err := CreateUser("username", "password", "password")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "username", user.Username)
	assert.NotNil(t, user.Password)
}

func TestChangePassword(t *testing.T) {
	// Given
	user, err := CreateUser("username", "password", "password")

	// When
	errPwd := user.ChangePassword("password", "newpassword", "newpassword")
	isValid, errValid := user.IsPasswordValid("newpassword")

	// Then
	assert.Nil(t, err)
	assert.Nil(t, errPwd)
	assert.Nil(t, errValid)
	assert.Equal(t, true, isValid)
}

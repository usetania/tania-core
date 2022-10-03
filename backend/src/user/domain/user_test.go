package domain_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	. "github.com/usetania/tania-core/src/user/domain"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) FindUserByUsername(username string) (UserServiceResult, error) {
	args := m.Called(username)

	return args.Get(0).(UserServiceResult), nil
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	// Given
	userServiceMock := new(UserServiceMock)
	userServiceMock.On("FindUserByUsername", "username").Return(UserServiceResult{})

	// When
	user, err := CreateUser(userServiceMock, "username", "password", "password")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "username", user.Username)
	assert.NotNil(t, user.Password)

	// Given
	userServiceMock2 := new(UserServiceMock)
	userUID, _ := uuid.NewV4()

	userServiceMock2.On("FindUserByUsername", "username").Return(UserServiceResult{
		UID:      userUID,
		Username: "username",
	})

	// When
	_, err = CreateUser(userServiceMock2, "username", "password", "password")

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, UserError{UserErrorUsernameExistsCode}, err)
}

func TestChangePassword(t *testing.T) {
	t.Parallel()
	// Given
	userServiceMock := new(UserServiceMock)
	userServiceMock.On("FindUserByUsername", "username").Return(UserServiceResult{})

	user, err := CreateUser(userServiceMock, "username", "password", "password")

	// When
	errPwd := user.ChangePassword("password", "newpassword", "newpassword")
	isValid, errValid := user.IsPasswordValid("newpassword")

	// Then
	assert.Nil(t, err)
	assert.Nil(t, errPwd)
	assert.Nil(t, errValid)
	assert.Equal(t, true, isValid)
}

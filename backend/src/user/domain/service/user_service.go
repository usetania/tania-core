package service

import (
	"errors"

	"github.com/usetania/tania-core/src/user/domain"
	"github.com/usetania/tania-core/src/user/query"
	"github.com/usetania/tania-core/src/user/storage"
)

type UserServiceImpl struct {
	UserReadQuery query.UserRead
}

func (s UserServiceImpl) FindUserByUsername(username string) (domain.UserServiceResult, error) {
	result := <-s.UserReadQuery.FindByUsername(username)

	if result.Error != nil {
		return domain.UserServiceResult{}, result.Error
	}

	user, ok := result.Result.(storage.UserRead)
	if !ok {
		return domain.UserServiceResult{}, errors.New("error type assertion")
	}

	return domain.UserServiceResult{
		UID:      user.UID,
		Username: user.Username,
	}, nil
}

package server

import (
	"errors"
	"log"

	"github.com/usetania/tania-core/src/user/domain"
	"github.com/usetania/tania-core/src/user/storage"
)

func (s *UserServer) SaveToUserReadModel(event interface{}) error {
	userRead := &storage.UserRead{}

	switch e := event.(type) {
	case domain.PasswordChanged:
		queryResult := <-s.UserReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			log.Println(queryResult.Error)
		}

		u, ok := queryResult.Result.(storage.UserRead)
		if !ok {
			log.Println(errors.New("internal server error. error type assertion"))
		}

		userRead = &u

		userRead.Password = e.NewPassword
		userRead.LastUpdated = e.DateChanged
	}

	err := <-s.UserReadRepo.Save(userRead)
	if err != nil {
		log.Println(err)
	}

	return nil
}

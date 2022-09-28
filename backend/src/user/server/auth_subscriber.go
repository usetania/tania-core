package server

import (
	"log"

	"github.com/usetania/tania-core/src/user/domain"
	"github.com/usetania/tania-core/src/user/storage"
)

func (s *AuthServer) SaveToUserReadModel(event interface{}) error {
	userRead := &storage.UserRead{}

	switch e := event.(type) {
	case domain.UserCreated:
		userRead.UID = e.UID
		userRead.Username = e.Username
		userRead.Password = e.Password
		userRead.CreatedDate = e.CreatedDate
		userRead.LastUpdated = e.LastUpdated
	}

	err := <-s.UserReadRepo.Save(userRead)
	if err != nil {
		log.Println(err)
	}

	return nil
}

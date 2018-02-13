package server

import (
	"fmt"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

func (s *FarmServer) SaveToFarmReadModel(event interface{}) error {
	fmt.Println("MASUK SUBSCRIBER")
	farmRead := &storage.FarmRead{}

	switch e := event.(type) {
	case domain.FarmCreated:
		farmRead.UID = e.UID
		farmRead.Name = e.Name
		farmRead.Type = e.Type
		farmRead.Latitude = e.Latitude
		farmRead.Longitude = e.Longitude
		farmRead.CountryCode = e.CountryCode
		farmRead.CityCode = e.CityCode
		farmRead.IsActive = e.IsActive
	}

	err := <-s.FarmReadRepo.Save(farmRead)
	if err != nil {
		return err
	}

	return nil
}

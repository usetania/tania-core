package server

import (
	"fmt"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

func (s *FarmServer) SaveToFarmReadModel(event interface{}) error {
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

func (s *FarmServer) SaveToReservoirReadModel(event interface{}) error {
	reservoirRead := &storage.ReservoirRead{}

	switch e := event.(type) {
	case domain.ReservoirCreated:
		fmt.Println("MASUK SINI DOONK")
		reservoirRead.UID = e.UID
		reservoirRead.Name = e.Name

		switch v := e.WaterSource.(type) {
		case domain.Bucket:
			reservoirRead.WaterSource = storage.WaterSource{
				Type:     v.Type(),
				Capacity: v.Capacity,
			}
		case domain.Tap:
			reservoirRead.WaterSource = storage.WaterSource{
				Type: v.Type(),
			}
		}

		reservoirRead.FarmUID = e.FarmUID
		reservoirRead.CreatedDate = e.CreatedDate
	}

	err := <-s.ReservoirReadRepo.Save(reservoirRead)
	if err != nil {
		return err
	}

	return nil
}

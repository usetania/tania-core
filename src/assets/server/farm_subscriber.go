package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/labstack/echo"
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
		farmRead.CreatedDate = e.CreatedDate
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
		queryResult := <-s.FarmReadQuery.FindByID(e.FarmUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		farm, ok := queryResult.Result.(storage.FarmRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

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

		reservoirRead.Farm = storage.ReservoirFarm{
			UID:  farm.UID,
			Name: farm.Name,
		}
		reservoirRead.CreatedDate = e.CreatedDate

	case domain.ReservoirNoteAdded:
		queryResult := <-s.ReservoirReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		r, ok := queryResult.Result.(storage.ReservoirRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		reservoirRead = &r

		reservoirRead.Notes = append(reservoirRead.Notes, storage.ReservoirNote{
			UID:         e.UID,
			Content:     e.Content,
			CreatedDate: e.CreatedDate,
		})

	case domain.ReservoirNoteRemoved:
		queryResult := <-s.ReservoirReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		r, ok := queryResult.Result.(query.ReservoirReadQueryResult)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		for _, v := range r.Notes {
			if v.UID != e.UID {
				reservoirRead.Notes = append(reservoirRead.Notes, storage.ReservoirNote{
					UID:         v.UID,
					Content:     v.Content,
					CreatedDate: v.CreatedDate,
				})
			}
		}

	}

	err := <-s.ReservoirReadRepo.Save(reservoirRead)
	if err != nil {
		return err
	}

	return nil
}

func (s *FarmServer) SaveToAreaReadModel(event interface{}) error {
	areaRead := &storage.AreaRead{}

	switch e := event.(type) {
	case domain.AreaCreated:
		queryResult := <-s.FarmReadQuery.FindByID(e.FarmUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		farm, ok := queryResult.Result.(storage.FarmRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		queryResult = <-s.ReservoirReadQuery.FindByID(e.ReservoirUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		reservoir, ok := queryResult.Result.(storage.ReservoirRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		areaRead.UID = e.UID
		areaRead.Name = e.Name
		areaRead.Type = storage.AreaType(e.Type)
		areaRead.Location = storage.AreaLocation(e.Location)
		areaRead.Size = storage.AreaSize(e.Size)
		areaRead.CreatedDate = e.CreatedDate
		areaRead.Farm = storage.AreaFarm{
			UID:  farm.UID,
			Name: farm.Name,
		}
		areaRead.Reservoir = storage.AreaReservoir{
			UID:  reservoir.UID,
			Name: reservoir.Name,
		}

	case domain.AreaPhotoAdded:
		queryResult := <-s.AreaReadQuery.FindByID(e.AreaUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		area, ok := queryResult.Result.(storage.AreaRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		areaRead = &area

		areaRead.Photo = storage.AreaPhoto{
			Filename: e.Filename,
			MimeType: e.MimeType,
			Size:     e.Size,
			Width:    e.Width,
			Height:   e.Height,
		}
	}

	err := <-s.AreaReadRepo.Save(areaRead)
	if err != nil {
		return err
	}

	return nil
}

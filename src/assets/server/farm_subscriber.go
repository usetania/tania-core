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

	case domain.ReservoirNoteAdded:
		queryResult := <-s.ReservoirReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		r, ok := queryResult.Result.(query.ReservoirReadQueryResult)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		mapReservoirReadToReservoirStorage(&r, reservoirRead)

		reservoirRead.Notes = append(reservoirRead.Notes, domain.ReservoirNote{
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

		mapReservoirReadToReservoirStorage(&r, reservoirRead)

		for _, v := range r.Notes {
			if v.UID != e.UID {
				reservoirRead.Notes = append(reservoirRead.Notes, domain.ReservoirNote{
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

func mapReservoirReadToReservoirStorage(
	rr *query.ReservoirReadQueryResult,
	rs *storage.ReservoirRead) *storage.ReservoirRead {

	rs.UID = rr.UID
	rs.Name = rr.Name
	rs.WaterSource = storage.WaterSource{
		Type:     rr.WaterSource.Type,
		Capacity: rr.WaterSource.Capacity,
	}
	rs.FarmUID = rr.FarmUID
	rs.CreatedDate = rr.CreatedDate

	return rs
}

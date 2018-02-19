package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/src/assets/domain"
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
		queryResult := <-s.ReservoirReadQuery.FindByID(e.ReservoirUID)
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
		queryResult := <-s.ReservoirReadQuery.FindByID(e.ReservoirUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		r, ok := queryResult.Result.(storage.ReservoirRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		reservoirRead = &r

		notes := []storage.ReservoirNote{}
		for _, v := range reservoirRead.Notes {
			if v.UID != e.UID {
				notes = append(notes, storage.ReservoirNote{
					UID:         v.UID,
					Content:     v.Content,
					CreatedDate: v.CreatedDate,
				})
			}
		}

		reservoirRead.Notes = notes

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

	case domain.AreaNoteAdded:
		queryResult := <-s.AreaReadQuery.FindByID(e.AreaUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		area, ok := queryResult.Result.(storage.AreaRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		areaRead = &area

		areaRead.Notes = append(areaRead.Notes, storage.AreaNote{
			UID:         e.UID,
			Content:     e.Content,
			CreatedDate: e.CreatedDated,
		})

	case domain.AreaNoteRemoved:
		queryResult := <-s.AreaReadQuery.FindByID(e.AreaUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		area, ok := queryResult.Result.(storage.AreaRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		areaRead = &area

		notes := []storage.AreaNote{}
		for _, v := range areaRead.Notes {
			if v.UID != e.UID {
				notes = append(notes, storage.AreaNote{
					UID:         v.UID,
					Content:     v.Content,
					CreatedDate: v.CreatedDate,
				})
			}
		}

		areaRead.Notes = notes
	}

	err := <-s.AreaReadRepo.Save(areaRead)
	if err != nil {
		return err
	}

	return nil
}

func (s *FarmServer) SaveToMaterialReadModel(event interface{}) error {
	materialRead := &storage.MaterialRead{}

	switch e := event.(type) {
	case domain.MaterialCreated:
		materialRead.UID = e.UID
		materialRead.Name = e.Name
		materialRead.PricePerUnit = e.PricePerUnit
		materialRead.Type = e.Type
		materialRead.Quantity = storage.MaterialQuantity(e.Quantity)
		materialRead.ExpirationDate = e.ExpirationDate
		materialRead.Notes = e.Notes
		materialRead.ProducedBy = e.ProducedBy
		materialRead.CreatedDate = e.CreatedDate

	case domain.MaterialNameChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.Name = e.Name

	case domain.MaterialPriceChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.PricePerUnit = e.Price

	case domain.MaterialQuantityChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.Quantity = storage.MaterialQuantity{
			Unit:  e.Quantity.Unit,
			Value: e.Quantity.Value,
		}

	case domain.MaterialTypeChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.Type = e.MaterialType

	case domain.MaterialExpirationDateChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.ExpirationDate = &e.ExpirationDate

	case domain.MaterialNotesChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.Notes = &e.Notes

	case domain.MaterialProducedByChanged:
		queryResult := <-s.MaterialReadQuery.FindByID(e.MaterialUID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		material, ok := queryResult.Result.(storage.MaterialRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		materialRead = &material

		materialRead.ProducedBy = &e.ProducedBy
	}

	err := <-s.MaterialReadRepo.Save(materialRead)
	if err != nil {
		return err
	}

	return nil
}

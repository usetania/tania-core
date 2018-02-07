package server

import (
	"net/http"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/labstack/echo"
)

func (s *GrowthServer) SaveToCropReadModel(event interface{}) error {
	cropRead := &storage.CropRead{}

	switch e := event.(type) {
	case domain.CropBatchCreated:
		cropRead.UID = e.UID
		cropRead.BatchID = e.BatchID
		cropRead.Status = e.Status.Code
		cropRead.Type = e.Type.Code
		cropRead.Container = storage.Container{
			Type:     e.ContainerType,
			Cell:     e.ContainerCell,
			Quantity: e.Quantity,
		}
		cropRead.Inventory = storage.Inventory{
			UID:       e.InventoryUID,
			Name:      e.VarietyName,
			PlantType: e.PlantType,
		}
		cropRead.InitialArea = storage.InitialArea{
			AreaUID:         e.InitialAreaUID,
			Name:            e.InitialAreaName,
			InitialQuantity: e.Quantity,
			CurrentQuantity: e.Quantity,
		}

		seeding := 0
		growing := 0
		if e.Type == domain.GetCropType(domain.CropTypeSeeding) {
			seeding += e.Quantity
		} else if e.Type == domain.GetCropType(domain.CropTypeGrowing) {
			growing += e.Quantity
		}

		cropRead.AreaStatus = storage.AreaStatus{
			Seeding: seeding,
			Growing: growing,
		}

		cropRead.FarmUID = e.FarmUID

	case domain.CropBatchWatered:
		queryResult := <-s.CropReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		cl, ok := queryResult.Result.(storage.CropRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		cropRead = &cl

		if cropRead.InitialArea.AreaUID == e.AreaUID {
			cropRead.InitialArea.LastWatered = &e.WateringDate
		}
	}

	err := <-s.CropReadRepo.Save(cropRead)
	if err != nil {
		return err
	}

	return nil
}

func (s *GrowthServer) SaveToCropActivityReadModel(event interface{}) error {
	cropActivity := &storage.CropActivity{}

	switch e := event.(type) {
	case domain.CropBatchCreated:
		cropActivity.UID = e.UID
		cropActivity.BatchID = e.BatchID
		cropActivity.ContainerType = e.ContainerType
		cropActivity.CreatedDate = e.CreatedDate
		cropActivity.ActivityType = storage.SeedActivity{
			AreaUID:  e.InitialAreaUID,
			AreaName: e.InitialAreaName,
			Quantity: e.Quantity,
		}
	case domain.CropBatchMoved:
		cropActivity.UID = e.UID
		cropActivity.BatchID = e.BatchID
		cropActivity.ContainerType = e.ContainerType
		cropActivity.CreatedDate = e.MovedDate
		cropActivity.ActivityType = storage.MoveActivity{
			SrcAreaUID:  e.SrcAreaUID,
			SrcAreaName: e.SrcAreaName,
			DstAreaUID:  e.DstAreaUID,
			DstAreaName: e.DstAreaName,
			Quantity:    e.Quantity,
		}
	}

	err := <-s.CropActivityRepo.Save(cropActivity)
	if err != nil {
		return err
	}

	return nil
}

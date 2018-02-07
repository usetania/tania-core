package server

import (
	"net/http"
	"time"

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

	case domain.CropBatchMoved:
		queryResult := <-s.CropReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		cr, ok := queryResult.Result.(storage.CropRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		cropRead = &cr

		if cropRead.InitialArea.AreaUID == e.SrcAreaUID {
			cropRead.InitialArea.CurrentQuantity -= e.Quantity
		}

		for i, v := range cropRead.MovedArea {
			if v.AreaUID == e.SrcAreaUID {
				cropRead.MovedArea[i].CurrentQuantity -= e.Quantity
			}
		}

		isDstExist := false
		for _, v := range cropRead.MovedArea {
			if v.AreaUID == e.DstAreaUID {
				isDstExist = true
			}
		}

		if isDstExist {
			for i, v := range cropRead.MovedArea {
				if v.AreaUID == e.DstAreaUID {
					cropRead.MovedArea[i].CurrentQuantity += e.Quantity
					cropRead.MovedArea[i].LastUpdated = e.MovedDate
				}
			}
		} else {
			cropRead.MovedArea = append(cropRead.MovedArea, storage.MovedArea{
				AreaUID:         e.DstAreaUID,
				Name:            e.DstAreaName,
				InitialQuantity: e.Quantity,
				CurrentQuantity: e.Quantity,
				CreatedDate:     e.MovedDate,
				LastUpdated:     e.MovedDate,
			})
		}

		if e.DstAreaType == "SEEDING" {
			cropRead.AreaStatus.Seeding += e.Quantity
		}
		if e.DstAreaType == "GROWING" {
			cropRead.AreaStatus.Growing += e.Quantity
		}
		if e.SrcAreaType == "SEEDING" {
			cropRead.AreaStatus.Seeding -= e.Quantity
		}
		if e.SrcAreaType == "GROWING" {
			cropRead.AreaStatus.Growing -= e.Quantity
		}

	case domain.CropBatchHarvested:
		queryResult := <-s.CropReadQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		cr, ok := queryResult.Result.(storage.CropRead)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		cropRead = &cr

		// If harvestType All, then empty the quantity in the area because it has been all harvested
		// Else if harvestType Partial, then we assume that the quantity of moved plant is 0
		harvestedQuantity := 0
		if e.HarvestType == domain.HarvestTypeAll {
			if cropRead.InitialArea.AreaUID == e.SrcAreaUID {
				harvestedQuantity = cropRead.InitialArea.CurrentQuantity
				cropRead.InitialArea.CurrentQuantity = 0
			}
			for i, v := range cropRead.MovedArea {
				if v.AreaUID == e.SrcAreaUID {
					harvestedQuantity = cropRead.MovedArea[i].CurrentQuantity
					cropRead.MovedArea[i].CurrentQuantity = 0
				}
			}
		}

		// Check source area existance. If already exist, then just update it
		isExist := false
		for i, v := range cropRead.HarvestedStorage {
			if v.SourceAreaUID == e.SrcAreaUID {
				cropRead.HarvestedStorage[i].Quantity += harvestedQuantity
				cropRead.HarvestedStorage[i].LastUpdated = e.HarvestDate
				isExist = true
			}
		}

		if !isExist {
			hs := storage.HarvestedStorage{
				Quantity:      harvestedQuantity,
				SourceAreaUID: e.SrcAreaUID,
				CreatedDate:   e.HarvestDate,
				LastUpdated:   e.HarvestDate,
			}
			cropRead.HarvestedStorage = append(cropRead.HarvestedStorage, hs)
		}

		// Calculate the produced harvest
		for i, v := range cropRead.HarvestedStorage {
			if v.SourceAreaUID == e.SrcAreaUID {
				cropRead.HarvestedStorage[i].ProducedGramQuantity += e.ProducedGramQuantity
			}
		}

		// Because Harvest should only be done in the GROWING area
		cropRead.AreaStatus.Growing -= e.Quantity

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
		cropActivity.CreatedDate = time.Now()
		cropActivity.ActivityType = storage.SeedActivity{
			AreaUID:     e.InitialAreaUID,
			AreaName:    e.InitialAreaName,
			Quantity:    e.Quantity,
			SeedingDate: e.CreatedDate,
		}
	case domain.CropBatchMoved:
		cropActivity.UID = e.UID
		cropActivity.BatchID = e.BatchID
		cropActivity.ContainerType = e.ContainerType
		cropActivity.CreatedDate = time.Now()
		cropActivity.ActivityType = storage.MoveActivity{
			SrcAreaUID:  e.SrcAreaUID,
			SrcAreaName: e.SrcAreaName,
			DstAreaUID:  e.DstAreaUID,
			DstAreaName: e.DstAreaName,
			Quantity:    e.Quantity,
			MovedDate:   e.MovedDate,
		}
	case domain.CropBatchHarvested:
		cropActivity.UID = e.UID
		cropActivity.BatchID = e.BatchID
		cropActivity.ContainerType = e.ContainerType
		cropActivity.CreatedDate = time.Now()
		cropActivity.ActivityType = storage.HarvestActivity{
			SrcAreaUID:           e.SrcAreaUID,
			SrcAreaName:          e.SrcAreaName,
			Quantity:             e.Quantity,
			ProducedGramQuantity: e.ProducedGramQuantity,
			HarvestType:          e.HarvestType,
			HarvestDate:          e.HarvestDate,
		}
	case domain.CropBatchWatered:
		cropActivity.UID = e.UID
		cropActivity.BatchID = e.BatchID
		cropActivity.ContainerType = e.ContainerType
		cropActivity.CreatedDate = time.Now()
		cropActivity.ActivityType = storage.WaterActivity{
			AreaUID:      e.AreaUID,
			AreaName:     e.AreaName,
			WateringDate: e.WateringDate,
		}
	}

	err := <-s.CropActivityRepo.Save(cropActivity)
	if err != nil {
		return err
	}

	return nil
}

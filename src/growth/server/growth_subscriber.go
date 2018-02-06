package server

import (
	"fmt"
	"net/http"

	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/labstack/echo"
)

func (s *GrowthServer) SaveToCropListReadModel(event interface{}) error {
	fmt.Println(event)
	cropList := &storage.CropList{}

	switch e := event.(type) {
	case domain.CropBatchCreated:
		cropList.UID = e.UID
		cropList.BatchID = e.BatchID
		cropList.VarietyName = e.VarietyName
		cropList.InventoryUID = e.InventoryUID
		cropList.InitialArea = storage.InitialArea{
			AreaUID: e.InitialAreaUID,
			Name:    e.InitialAreaName,
			InitialQuantity: storage.Container{
				Type:     e.ContainerType,
				Quantity: e.Quantity,
				Cell:     e.ContainerCell,
			},
			CurrentQuantity: storage.Container{
				Type:     e.ContainerType,
				Quantity: e.Quantity,
				Cell:     e.ContainerCell,
			},
		}

		seeding := 0
		growing := 0
		if e.Type == domain.GetCropType(domain.CropTypeSeeding) {
			seeding += e.Quantity
		} else if e.Type == domain.GetCropType(domain.CropTypeGrowing) {
			growing += e.Quantity
		}

		cropList.AreaStatus = storage.AreaStatus{
			Seeding: seeding,
			Growing: growing,
		}

		cropList.CreatedDate = e.CreatedDate
		cropList.FarmUID = e.FarmUID

	case domain.CropBatchWatered:
		fmt.Println("CROP BATCH WATERED")
		queryResult := <-s.CropListQuery.FindByID(e.UID)
		if queryResult.Error != nil {
			return queryResult.Error
		}

		cl, ok := queryResult.Result.(storage.CropList)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		cropList = &cl

		if cropList.InitialArea.AreaUID == e.AreaUID {
			cropList.InitialArea.LastWatered = &e.WateringDate
		}
	}
	fmt.Println(cropList)
	err := <-s.CropListRepo.Save(cropList)
	if err != nil {
		return err
	}

	return nil
}

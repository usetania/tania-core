package server

import (
	"github.com/Tanibox/tania-server/src/growth/domain"
	growthdomain "github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
)

func (s *GrowthServer) SaveToCropListReadModel(event growthdomain.CropBatchCreated) error {
	cropList := &storage.CropList{}
	cropList.UID = event.UID
	cropList.BatchID = event.BatchID
	cropList.VarietyName = event.VarietyName
	cropList.InventoryUID = event.InventoryUID
	cropList.InitialArea = storage.InitialArea{
		AreaUID: event.InitialAreaUID,
		Name:    event.InitialAreaName,
		InitialQuantity: storage.Container{
			Type:     event.ContainerType,
			Quantity: event.Quantity,
			Cell:     event.ContainerCell,
		},
		CurrentQuantity: storage.Container{
			Type:     event.ContainerType,
			Quantity: event.Quantity,
			Cell:     event.ContainerCell,
		},
	}

	seeding := 0
	growing := 0
	if event.Type == domain.GetCropType(domain.CropTypeSeeding) {
		seeding += event.Quantity
	} else if event.Type == domain.GetCropType(domain.CropTypeGrowing) {
		growing += event.Quantity
	}

	cropList.AreaStatus = storage.AreaStatus{
		Seeding: seeding,
		Growing: growing,
	}

	cropList.CreatedDate = event.CreatedDate
	cropList.FarmUID = event.FarmUID

	err := <-s.CropListRepo.Save(cropList)
	if err != nil {
		return err
	}

	return nil
}

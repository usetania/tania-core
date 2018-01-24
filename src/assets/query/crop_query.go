package query

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	growthdomain "github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropQuery interface {
	FindAllCropByArea(areaUID uuid.UUID) <-chan QueryResult
	CountCropsByArea(areaUID uuid.UUID) <-chan QueryResult
}

type CropQueryInMemory struct {
	Storage *storage.CropStorage
}

func NewCropQueryInMemory(s *storage.CropStorage) CropQuery {
	return CropQueryInMemory{Storage: s}
}

func (q CropQueryInMemory) CountCropsByArea(areaUID uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		totalCropBatch := 0
		totalPlant := 0
		for _, val := range q.Storage.CropMap {
			if val.InitialArea.AreaUID == areaUID {
				totalCropBatch++
				totalPlant += val.Container.Quantity
			}

			for _, v := range val.MovedArea {
				if v.AreaUID == areaUID {
					totalCropBatch++
					totalPlant += val.Container.Quantity
				}
			}
		}

		result <- QueryResult{Result: domain.CountAreaCrop{
			PlantQuantity:  totalPlant,
			TotalCropBatch: totalCropBatch,
		}}

		close(result)
	}()

	return result
}

func (q CropQueryInMemory) FindAllCropByArea(areaUID uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		crops := []domain.AreaCrop{}
		for _, val := range q.Storage.CropMap {
			if val.InitialArea.AreaUID == areaUID {
				containerType := domain.CropContainerType{}
				switch v := val.Container.Type.(type) {
				case growthdomain.Tray:
					containerType.Code = v.Code()
					containerType.Cell = v.Cell
				case growthdomain.Pot:
					containerType.Code = v.Code()
				}

				crops = append(crops, domain.AreaCrop{
					CropUID: val.UID,
					BatchID: val.BatchID,
					InitialArea: domain.InitialArea{
						AreaUID: val.InitialArea.AreaUID,
					},
					MovingDate:  val.CreatedDate,
					CreatedDate: val.CreatedDate,
					Inventory: domain.CropInventory{
						UID: val.InventoryUID,
					},
					Container: domain.CropContainer{
						Quantity: val.Container.Quantity,
						Type:     containerType,
					},
				})
			}

			for _, v := range val.MovedArea {
				if v.AreaUID == areaUID {
					containerType := domain.CropContainerType{}
					switch v := val.Container.Type.(type) {
					case growthdomain.Tray:
						containerType.Code = v.Code()
						containerType.Cell = v.Cell
					case growthdomain.Pot:
						containerType.Code = v.Code()
					}

					crops = append(crops, domain.AreaCrop{
						CropUID: val.UID,
						BatchID: val.BatchID,
						InitialArea: domain.InitialArea{
							AreaUID: v.SourceAreaUID,
						},
						MovingDate:  val.CreatedDate,
						CreatedDate: val.CreatedDate,
						Inventory: domain.CropInventory{
							UID: val.InventoryUID,
						},
						Container: domain.CropContainer{
							Quantity: val.Container.Quantity,
							Type:     containerType,
						},
					})
				}
			}
		}

		result <- QueryResult{Result: crops}

		close(result)
	}()

	return result
}

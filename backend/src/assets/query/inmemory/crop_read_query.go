package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropReadQueryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropReadQueryInMemory(s *storage.CropReadStorage) query.CropReadQuery {
	return CropReadQueryInMemory{Storage: s}
}

func (q CropReadQueryInMemory) CountCropsByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		totalCropBatch := 0
		totalPlant := 0
		for _, val := range q.Storage.CropReadMap {
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

		result <- query.QueryResult{Result: query.CountAreaCropQueryResult{
			PlantQuantity:  totalPlant,
			TotalCropBatch: totalCropBatch,
		}}

		close(result)
	}()

	return result
}

func (q CropReadQueryInMemory) FindAllCropByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		crops := []query.AreaCropQueryResult{}
		for _, val := range q.Storage.CropReadMap {
			if val.InitialArea.AreaUID == areaUID {
				crops = append(crops, query.AreaCropQueryResult{
					CropUID: val.UID,
					BatchID: val.BatchID,
					InitialArea: query.InitialArea{
						AreaUID: val.InitialArea.AreaUID,
						Name:    "",
					},
					MovingDate:  val.InitialArea.CreatedDate,
					CreatedDate: val.InitialArea.CreatedDate,
					Inventory: query.Inventory{
						UID: val.Inventory.UID,
					},
					Container: query.Container{
						Quantity: val.Container.Quantity,
						Type: query.ContainerType{
							Code: val.Container.Type,
							Cell: val.Container.Cell,
						},
					},
				})
			}

			for _, v := range val.MovedArea {
				if v.AreaUID == areaUID {
					crops = append(crops, query.AreaCropQueryResult{
						CropUID: val.UID,
						BatchID: val.BatchID,
						InitialArea: query.InitialArea{
							AreaUID: val.InitialArea.AreaUID,
						},
						MovingDate:  val.InitialArea.CreatedDate,
						CreatedDate: val.InitialArea.CreatedDate,
						Inventory: query.Inventory{
							UID: val.Inventory.UID,
						},
						Container: query.Container{
							Quantity: val.Container.Quantity,
							Type: query.ContainerType{
								Code: val.Container.Type,
								Cell: val.Container.Cell,
							},
						},
					})
				}
			}
		}

		result <- query.QueryResult{Result: crops}

		close(result)
	}()

	return result
}

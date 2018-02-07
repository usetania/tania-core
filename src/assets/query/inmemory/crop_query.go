package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/query"
	growthdomain "github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropQueryInMemory struct {
	Storage *storage.CropStorage
}

func NewCropQueryInMemory(s *storage.CropStorage) query.CropQuery {
	return CropQueryInMemory{Storage: s}
}

func (q CropQueryInMemory) CountCropsByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

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

		result <- query.QueryResult{Result: query.CountAreaCropQueryResult{
			PlantQuantity:  totalPlant,
			TotalCropBatch: totalCropBatch,
		}}

		close(result)
	}()

	return result
}

func (q CropQueryInMemory) FindAllCropByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		crops := []query.AreaCropQueryResult{}
		for _, val := range q.Storage.CropMap {
			if val.InitialArea.AreaUID == areaUID {

				containerTypeCode := ""
				containerTypeCell := 0
				switch v := val.Container.Type.(type) {
				case growthdomain.Tray:
					containerTypeCode = v.Code()
					containerTypeCell = v.Cell
				case growthdomain.Pot:
					containerTypeCode = v.Code()
				}

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
						UID: val.InventoryUID,
					},
					Container: query.Container{
						Quantity: val.Container.Quantity,
						Type: query.ContainerType{
							Code: containerTypeCode,
							Cell: containerTypeCell,
						},
					},
				})
			}

			for _, v := range val.MovedArea {
				if v.AreaUID == areaUID {

					containerTypeCode := ""
					containerTypeCell := 0
					switch v := val.Container.Type.(type) {
					case growthdomain.Tray:
						containerTypeCode = v.Code()
						containerTypeCell = v.Cell
					case growthdomain.Pot:
						containerTypeCode = v.Code()
					}

					crops = append(crops, query.AreaCropQueryResult{
						CropUID: val.UID,
						BatchID: val.BatchID,
						InitialArea: query.InitialArea{
							AreaUID: v.SourceAreaUID,
						},
						MovingDate:  val.InitialArea.CreatedDate,
						CreatedDate: val.InitialArea.CreatedDate,
						Inventory: query.Inventory{
							UID: val.InventoryUID,
						},
						Container: query.Container{
							Quantity: val.Container.Quantity,
							Type: query.ContainerType{
								Code: containerTypeCode,
								Cell: containerTypeCell,
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

package inmemory

import (
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropReadQueryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropReadQueryInMemory(s *storage.CropReadStorage) query.CropReadQuery {
	return CropReadQueryInMemory{Storage: s}
}

func (s CropReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
			if val.UID == uid {
				crop = val
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindByBatchID(batchID string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
			if val.BatchID == batchID {
				crop = val
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		cropRead := []storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				cropRead = append(cropRead, val)
			}
		}

		result <- query.QueryResult{Result: cropRead}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crops := []query.CropAreaByAreaQueryResult{}
		for _, val := range s.Storage.CropReadMap {
			if val.InitialArea.AreaUID == areaUID {
				crops = append(crops, query.CropAreaByAreaQueryResult{
					UID:         val.UID,
					BatchID:     val.BatchID,
					CreatedDate: val.InitialArea.CreatedDate,
					Area: query.Area{
						UID:             val.InitialArea.AreaUID,
						Name:            val.InitialArea.Name,
						InitialQuantity: val.InitialArea.InitialQuantity,
						CurrentQuantity: val.InitialArea.CurrentQuantity,
						InitialArea: query.InitialArea{
							UID:         val.InitialArea.AreaUID,
							Name:        val.InitialArea.Name,
							CreatedDate: val.InitialArea.CreatedDate,
						},
						LastWatered: val.InitialArea.LastWatered,
						MovingDate:  val.InitialArea.CreatedDate,
					},
					Container: query.Container{
						Type:     val.Container.Type,
						Cell:     val.Container.Cell,
						Quantity: val.Container.Quantity,
					},
					Inventory: query.Inventory{
						UID:       val.Inventory.UID,
						Name:      val.Inventory.Name,
						PlantType: val.Inventory.PlantType,
					},
				})
			}

			for i := range val.MovedArea {
				if val.MovedArea[i].AreaUID == areaUID {
					crops = append(crops, query.CropAreaByAreaQueryResult{
						UID:         val.UID,
						BatchID:     val.BatchID,
						CreatedDate: val.MovedArea[i].CreatedDate,
						Area: query.Area{
							UID:             val.MovedArea[i].AreaUID,
							Name:            val.MovedArea[i].Name,
							InitialQuantity: val.MovedArea[i].InitialQuantity,
							CurrentQuantity: val.MovedArea[i].CurrentQuantity,
							InitialArea: query.InitialArea{
								UID:         val.InitialArea.AreaUID,
								Name:        val.InitialArea.Name,
								CreatedDate: val.InitialArea.CreatedDate,
							},
							LastWatered: val.MovedArea[i].LastWatered,
							MovingDate:  val.MovedArea[i].CreatedDate,
						},
						Container: query.Container{
							Type:     val.Container.Type,
							Cell:     val.Container.Cell,
							Quantity: val.Container.Quantity,
						},
						Inventory: query.Inventory{
							UID:       val.Inventory.UID,
							Name:      val.Inventory.Name,
							PlantType: val.Inventory.PlantType,
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

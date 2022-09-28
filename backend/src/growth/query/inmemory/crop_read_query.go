package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropReadQueryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropReadQueryInMemory(s *storage.CropReadStorage) query.CropReadQuery {
	return CropReadQueryInMemory{Storage: s}
}

func (s CropReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropRead{}

		for _, val := range s.Storage.CropReadMap {
			if val.UID == uid {
				crop = val
			}
		}

		result <- query.Result{Result: crop}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindByBatchID(batchID string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropRead{}

		for _, val := range s.Storage.CropReadMap {
			if val.BatchID == batchID {
				crop = val
			}
		}

		result <- query.Result{Result: crop}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsByFarm(farmUID uuid.UUID, _ string, _, _ int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		cropRead := []storage.CropRead{}

		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				// Check all the current quantity
				// It should not be zero,
				// because if all zero then it will show up in the Archieves instead
				initialEmpty := false
				allMovedEmpty := []bool{}

				if val.InitialArea.CurrentQuantity <= 0 {
					initialEmpty = true
				}

				for _, v := range val.MovedArea {
					if v.CurrentQuantity <= 0 {
						allMovedEmpty = append(allMovedEmpty, true)
					}
				}

				movedEmpty := false

				for _, v := range allMovedEmpty {
					if v {
						movedEmpty = true
					}
				}

				if !initialEmpty || !movedEmpty {
					cropRead = append(cropRead, val)
				}
			}
		}

		result <- query.Result{Result: cropRead}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) CountAllCropsByFarm(_ uuid.UUID, _ string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		total := len(s.Storage.CropReadMap)
		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsArchives(farmUID uuid.UUID, _, _ int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		archives := []storage.CropRead{}

		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				// A crop's current quantity which have zero value should go to archives
				initialEmpty := true
				allMovedEmpty := []bool{}

				if val.InitialArea.CurrentQuantity > 0 {
					initialEmpty = false
				}

				for _, v := range val.MovedArea {
					if v.CurrentQuantity > 0 {
						allMovedEmpty = append(allMovedEmpty, false)
					}
				}

				if initialEmpty {
					movedEmpty := true

					for _, v := range allMovedEmpty {
						if !v {
							movedEmpty = false
						}
					}

					if movedEmpty {
						archives = append(archives, val)
					}
				}
			}
		}

		result <- query.Result{Result: archives}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) CountAllArchivedCropsByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		total := 0

		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				// A crop's current quantity which have zero value should go to archives
				initialEmpty := true
				allMovedEmpty := []bool{}

				if val.InitialArea.CurrentQuantity > 0 {
					initialEmpty = false
				}

				for _, v := range val.MovedArea {
					if v.CurrentQuantity > 0 {
						allMovedEmpty = append(allMovedEmpty, false)
					}
				}

				if initialEmpty {
					movedEmpty := true

					for _, v := range allMovedEmpty {
						if !v {
							movedEmpty = false
						}
					}

					if movedEmpty {
						total++
					}
				}
			}
		}

		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsByArea(areaUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

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

		result <- query.Result{Result: crops}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindCropsInformation(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		cropInf := query.CropInformationQueryResult{}
		harvestProduced := float32(0)
		plantType := make(map[string]bool)
		totalPlantVariety := 0

		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				for _, val2 := range val.HarvestedStorage {
					harvestProduced += val2.ProducedGramQuantity
				}

				if _, ok := plantType[val.Inventory.Name]; !ok {
					totalPlantVariety++

					plantType[val.Inventory.Name] = true
				}
			}
		}

		cropInf.TotalHarvestProduced = harvestProduced
		cropInf.TotalPlantVariety = totalPlantVariety

		result <- query.Result{Result: cropInf}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) CountTotalBatch(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		varQty := []query.CountTotalBatchQueryResult{}
		varietyName := make(map[string]int)

		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				if _, ok := varietyName[val.Inventory.Name]; !ok {
					varietyName[val.Inventory.Name]++
				}
			}
		}

		for i, v := range varietyName {
			varQty = append(varQty, query.CountTotalBatchQueryResult{
				VarietyName: i,
				TotalBatch:  v,
			})
		}

		result <- query.Result{Result: varQty}

		close(result)
	}()

	return result
}

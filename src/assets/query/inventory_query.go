package query

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type InventoryMaterialQuery interface {
	FindAllInventoryByPlantType(plantType domain.PlantType) <-chan QueryResult
	FindInventoryByPlantTypeAndVariety(plantType domain.PlantType, variety string) <-chan QueryResult
}

type InventoryMaterialQueryInMemory struct {
	Storage *storage.InventoryMaterialStorage
}

func NewInventoryMaterialQueryInMemory(s *storage.InventoryMaterialStorage) InventoryMaterialQuery {
	return &InventoryMaterialQueryInMemory{Storage: s}
}

func (q *InventoryMaterialQueryInMemory) FindAllInventoryByPlantType(plantType domain.PlantType) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		inventories := []domain.InventoryMaterial{}
		for _, val := range q.Storage.InventoryMaterialMap {
			if val.PlantType == plantType {
				inventories = append(inventories, val)
			}
		}

		result <- QueryResult{Result: inventories}

		close(result)
	}()

	return result
}

func (q *InventoryMaterialQueryInMemory) FindInventoryByPlantTypeAndVariety(plantType domain.PlantType, variety string) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		inventory := domain.InventoryMaterial{}
		for _, val := range q.Storage.InventoryMaterialMap {
			if val.PlantType == plantType && val.Variety == variety {
				inventory = val
			}
		}

		result <- QueryResult{Result: inventory}

		close(result)
	}()

	return result
}

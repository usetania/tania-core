package query

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type InventoryMaterialQuery interface {
	FindAllInventoryByPlantType(plantType domain.PlantType) <-chan QueryResult
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

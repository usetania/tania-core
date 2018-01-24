package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type InventoryMaterialQueryInMemory struct {
	Storage *storage.InventoryMaterialStorage
}

func NewInventoryMaterialQueryInMemory(s *storage.InventoryMaterialStorage) query.InventoryMaterialQuery {
	return &InventoryMaterialQueryInMemory{Storage: s}
}

func (q *InventoryMaterialQueryInMemory) FindAllInventoryByPlantType(plantType domain.PlantType) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		inventories := []domain.InventoryMaterial{}
		for _, val := range q.Storage.InventoryMaterialMap {
			if val.PlantType == plantType {
				inventories = append(inventories, val)
			}
		}

		result <- query.QueryResult{Result: inventories}

		close(result)
	}()

	return result
}

func (q *InventoryMaterialQueryInMemory) FindInventoryByPlantTypeAndVariety(plantType domain.PlantType, variety string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		inventory := domain.InventoryMaterial{}
		for _, val := range q.Storage.InventoryMaterialMap {
			if val.PlantType == plantType && val.Variety == variety {
				inventory = val
			}
		}

		result <- query.QueryResult{Result: inventory}

		close(result)
	}()

	return result
}

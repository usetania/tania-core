package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type InventoryMaterialQueryInMemory struct {
	Storage *storage.MaterialStorage
}

func NewInventoryMaterialQueryInMemory(s *storage.MaterialStorage) query.InventoryMaterialQuery {
	return InventoryMaterialQueryInMemory{Storage: s}
}

func (s InventoryMaterialQueryInMemory) FindByID(inventoryUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		ci := query.CropInventoryQueryResult{}
		for _, val := range s.Storage.MaterialMap {
			if val.UID == inventoryUID {
				ci.UID = val.UID
				// ci.PlantTypeCode = val.PlantType.Code()
				// ci.Variety = val.Variety
			}
		}

		result <- query.QueryResult{Result: ci}

		close(result)
	}()

	return result
}

func (q InventoryMaterialQueryInMemory) FindInventoryByPlantTypeCodeAndVariety(plantTypeCode string, variety string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		ci := query.CropInventoryQueryResult{}
		for _, val := range q.Storage.MaterialMap {
			// if val.PlantType.Code() == plantTypeCode && val.Variety == variety {
			ci.UID = val.UID
			// ci.PlantTypeCode = val.PlantType.Code()
			// ci.Variety = val.Variety
			// }
		}

		result <- query.QueryResult{Result: ci}

		close(result)
	}()

	return result
}

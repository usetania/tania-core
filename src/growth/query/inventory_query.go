package query

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type InventoryMaterialQuery interface {
	FindByID(inventoryUID uuid.UUID) <-chan QueryResult
}

type InventoryMaterialQueryInMemory struct {
	Storage *storage.InventoryMaterialStorage
}

type CropInventory struct {
	UID       uuid.UUID
	PlantType string
	Variety   string
}

func NewInventoryMaterialQueryInMemory(s *storage.InventoryMaterialStorage) InventoryMaterialQuery {
	return InventoryMaterialQueryInMemory{Storage: s}
}

func (s InventoryMaterialQueryInMemory) FindByID(inventoryUID uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		ci := CropInventory{}
		for _, val := range s.Storage.InventoryMaterialMap {
			if val.UID == inventoryUID {
				ci.UID = val.UID
				ci.PlantType = val.PlantType.Code()
				ci.Variety = val.Variety
			}
		}

		result <- QueryResult{Result: ci}

		close(result)
	}()

	return result
}

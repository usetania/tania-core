package inmemory

import (
	"github.com/Tanibox/tania-server/src/growth/storage"
	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type CropQueryInMemory struct {
	Storage *storage.CropStorage
}

func NewCropQueryInMemory(s *storage.CropStorage) query.CropQuery {
	return CropQueryInMemory{Storage: s}
}

func (s CropQueryInMemory) FindCropByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := query.CropcropQueryResult{}
		for _, val := range s.Storage.cropMap {
			if val.UID == uid {
				crop.UID = uid
				crop.BatchID = val.BatchID
				crop.Status = val.Status
				crop.Type = val.Type
				crop.Container.Quantity = val.Container.Quantity
				crop.Container.Type = val.Container.Type
				crop.InventoryUID = val.Inventory.UID
				crop.FarmUID = val.Farm.UID
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

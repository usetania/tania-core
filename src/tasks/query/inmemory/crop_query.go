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

		crop := query.TaskCropQueryResult{}
		for _, val := range s.Storage.CropMap {
			if val.UID == uid {
				crop.UID = uid
				crop.BatchID = val.BatchID
				crop.Status = val.Status.Code
				crop.Type = val.Type.Code
				crop.Container.Quantity = val.Container.Quantity
				crop.Container.Type = val.Container.Type.Code()
				crop.InventoryUID = val.InventoryUID
				crop.FarmUID = val.FarmUID
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

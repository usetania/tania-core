package inmemory

import (
	"github.com/Tanibox/tania-core/src/growth/storage"
	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type CropQueryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropQueryInMemory(s *storage.CropReadStorage) query.CropQuery {
	return CropQueryInMemory{Storage: s}
}

func (s CropQueryInMemory) FindCropByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := query.TaskCropQueryResult{}

		for _, val := range s.Storage.CropReadMap {
			if val.UID == uid {
				crop.UID = uid
				crop.BatchID = val.BatchID
			}
		}
		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

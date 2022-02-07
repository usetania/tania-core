package inmemory

import (
	"github.com/Tanibox/tania-core/src/growth/storage"
	"github.com/Tanibox/tania-core/src/tasks/query"
	"github.com/gofrs/uuid"
)

type CropQueryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropQueryInMemory(s *storage.CropReadStorage) query.Crop {
	return CropQueryInMemory{Storage: s}
}

func (s CropQueryInMemory) FindCropByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := query.TaskCropResult{}

		for _, val := range s.Storage.CropReadMap {
			if val.UID == uid {
				crop.UID = uid
				crop.BatchID = val.BatchID
			}
		}
		result <- query.Result{Result: crop}

		close(result)
	}()

	return result
}

package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/storage"
	"github.com/usetania/tania-core/src/tasks/query"
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

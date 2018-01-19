package query

import (
	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
)

type CropQuery interface {
	FindByBatchID(batchID string) <-chan QueryResult
}

type CropQueryInMemory struct {
	Storage *storage.CropStorage
}

func NewCropQueryInMemory(s *storage.CropStorage) CropQuery {
	return CropQueryInMemory{Storage: s}
}

func (s CropQueryInMemory) FindByBatchID(batchID string) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := domain.Crop{}
		for _, val := range s.Storage.CropMap {
			if val.BatchID == batchID {
				crop = val
			}
		}

		result <- QueryResult{Result: crop}

		close(result)
	}()

	return result
}

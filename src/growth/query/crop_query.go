package query

import (
	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropQuery interface {
	FindByBatchID(batchID string) <-chan QueryResult
	FindAllCropsByFarm(farmUID uuid.UUID) <-chan QueryResult
	FindAllCropsByArea(areaUID uuid.UUID) <-chan QueryResult
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

func (s CropQueryInMemory) FindAllCropsByFarm(farmUID uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crops := []domain.Crop{}
		for _, val := range s.Storage.CropMap {
			if val.FarmUID == farmUID {
				crops = append(crops, val)
			}
		}

		result <- QueryResult{Result: crops}

		close(result)
	}()

	return result
}

func (s CropQueryInMemory) FindAllCropsByArea(areaUID uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crops := []domain.Crop{}
		for _, val := range s.Storage.CropMap {
			if val.InitialArea.AreaUID == areaUID {
				crops = append(crops, val)
			}

			for _, val2 := range val.MovedArea {
				if val2.AreaUID == areaUID {
					crops = append(crops, val)
				}
			}

		}

		result <- QueryResult{Result: crops}

		close(result)
	}()

	return result
}

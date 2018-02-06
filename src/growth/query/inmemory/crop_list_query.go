package inmemory

import (
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropListQueryInMemory struct {
	Storage *storage.CropListStorage
}

func NewCropListQueryInMemory(s *storage.CropListStorage) query.CropListQuery {
	return CropListQueryInMemory{Storage: s}
}

func (s CropListQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropList{}
		for _, val := range s.Storage.CropListMap {
			if val.UID == uid {
				crop = val
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

func (s CropListQueryInMemory) FindByBatchID(batchID string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropList{}
		for _, val := range s.Storage.CropListMap {
			if val.BatchID == batchID {
				crop = val
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

func (s CropListQueryInMemory) FindAllCropsByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		cropList := []storage.CropList{}
		for _, val := range s.Storage.CropListMap {
			if val.FarmUID == farmUID {
				cropList = append(cropList, val)
			}
		}

		result <- query.QueryResult{Result: cropList}

		close(result)
	}()

	return result
}

func (s CropListQueryInMemory) FindAllCropsByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crops := []storage.CropList{}
		for _, val := range s.Storage.CropListMap {
			if val.InitialArea.AreaUID == areaUID {
				crops = append(crops, val)
			}

			for _, val2 := range val.MovedArea {
				if val2.AreaUID == areaUID {
					crops = append(crops, val)
				}
			}

		}

		result <- query.QueryResult{Result: crops}

		close(result)
	}()

	return result
}

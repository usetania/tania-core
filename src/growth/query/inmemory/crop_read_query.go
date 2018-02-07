package inmemory

import (
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropReadQueryInMemory struct {
	Storage *storage.CropReadStorage
}

func NewCropListQueryInMemory(s *storage.CropReadStorage) query.CropListQuery {
	return CropReadQueryInMemory{Storage: s}
}

func (s CropReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
			if val.UID == uid {
				crop = val
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindByBatchID(batchID string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crop := storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
			if val.BatchID == batchID {
				crop = val
			}
		}

		result <- query.QueryResult{Result: crop}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		cropList := []storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
			if val.FarmUID == farmUID {
				cropList = append(cropList, val)
			}
		}

		result <- query.QueryResult{Result: cropList}

		close(result)
	}()

	return result
}

func (s CropReadQueryInMemory) FindAllCropsByArea(areaUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crops := []storage.CropRead{}
		for _, val := range s.Storage.CropReadMap {
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

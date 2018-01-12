package query

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type CropQuery interface {
	FindAllCropsByArea(area domain.Area) <-chan QueryResult
}

type CropQueryInMemory struct {
	Storage *storage.CropStorage
}

func NewCropQueryInMemory(s *storage.CropStorage) CropQuery {
	return CropQueryInMemory{Storage: s}
}

func (s CropQueryInMemory) FindAllCropsByArea(area domain.Area) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		crops := []domain.Crop{}
		for _, val := range s.Storage.CropMap {
			for _, currArea := range val.CurrentAreas {
				if currArea.UID == area.UID {
					crops = append(crops, val)
				}
			}
		}

		result <- QueryResult{Result: crops}

		close(result)
	}()

	return result
}

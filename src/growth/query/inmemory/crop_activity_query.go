package inmemory

import (
	"github.com/Tanibox/tania-server/src/growth/query"
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type CropActivityQueryInMemory struct {
	Storage *storage.CropActivityStorage
}

func NewCropActivityQueryInMemory(s *storage.CropActivityStorage) query.CropActivityQuery {
	return CropActivityQueryInMemory{Storage: s}
}

func (s CropActivityQueryInMemory) FindAllByCropID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		activities := []storage.CropActivity{}
		for _, val := range s.Storage.CropActivityMap {
			if val.UID == uid {
				activities = append(activities, val)
			}
		}

		result <- query.QueryResult{Result: activities}

		close(result)
	}()

	return result
}

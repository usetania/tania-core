package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type FarmQueryInMemory struct {
	Storage *storage.FarmStorage
}

func NewFarmQueryInMemory(s *storage.FarmStorage) query.FarmQuery {
	return FarmQueryInMemory{Storage: s}
}

func (s FarmQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farm := query.CropFarmQueryResult{}
		for _, val := range s.Storage.FarmMap {
			if val.UID == uid {
				farm.UID = uid
				farm.Name = val.Name
			}
		}

		result <- query.QueryResult{Result: farm}

		close(result)
	}()

	return result
}

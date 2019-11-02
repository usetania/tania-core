package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmReadQueryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadQueryInMemory(s *storage.FarmReadStorage) query.FarmReadQuery {
	return FarmReadQueryInMemory{Storage: s}
}

func (s FarmReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farm := storage.FarmRead{}
		for _, val := range s.Storage.FarmReadMap {
			if val.UID == uid {
				farm = val
			}
		}

		result <- query.QueryResult{Result: farm}

		close(result)
	}()

	return result
}

func (s FarmReadQueryInMemory) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farms := []storage.FarmRead{}
		for _, val := range s.Storage.FarmReadMap {
			farms = append(farms, val)
		}

		result <- query.QueryResult{Result: farms}

		close(result)
	}()

	return result
}

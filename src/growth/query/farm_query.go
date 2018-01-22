package query

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/growth/domain"
	uuid "github.com/satori/go.uuid"
)

type FarmQuery interface {
	FindByID(farmUID uuid.UUID) <-chan QueryResult
}

type FarmQueryInMemory struct {
	Storage *storage.FarmStorage
}

func NewFarmQueryInMemory(s *storage.FarmStorage) FarmQuery {
	return FarmQueryInMemory{Storage: s}
}

func (s FarmQueryInMemory) FindByID(uid uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farm := domain.CropFarm{}
		for _, val := range s.Storage.FarmMap {
			if val.UID == uid {
				farm.UID = uid
				farm.Name = val.Name
			}
		}

		result <- QueryResult{Result: farm}

		close(result)
	}()

	return result
}

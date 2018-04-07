package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type ReservoirReadQueryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirReadQueryInMemory(s *storage.ReservoirReadStorage) query.ReservoirReadQuery {
	return ReservoirReadQueryInMemory{Storage: s}
}

func (s ReservoirReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		reservoir := storage.ReservoirRead{}
		for _, val := range s.Storage.ReservoirReadMap {
			if val.UID == uid {
				reservoir = val
			}
		}

		result <- query.QueryResult{Result: reservoir}

		close(result)
	}()

	return result
}

func (s ReservoirReadQueryInMemory) FindAllByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		reservoirs := []storage.ReservoirRead{}
		for _, val := range s.Storage.ReservoirReadMap {
			if val.Farm.UID == farmUID {
				reservoirs = append(reservoirs, val)
			}
		}

		result <- query.QueryResult{Result: reservoirs}

		close(result)
	}()

	return result
}

package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirReadQueryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirReadQueryInMemory(s *storage.ReservoirReadStorage) query.ReservoirRead {
	return ReservoirReadQueryInMemory{Storage: s}
}

func (s ReservoirReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		reservoir := storage.ReservoirRead{}

		for _, val := range s.Storage.ReservoirReadMap {
			if val.UID == uid {
				reservoir = val
			}
		}

		result <- query.Result{Result: reservoir}

		close(result)
	}()

	return result
}

func (s ReservoirReadQueryInMemory) FindAllByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		reservoirs := []storage.ReservoirRead{}

		for _, val := range s.Storage.ReservoirReadMap {
			if val.Farm.UID == farmUID {
				reservoirs = append(reservoirs, val)
			}
		}

		result <- query.Result{Result: reservoirs}

		close(result)
	}()

	return result
}

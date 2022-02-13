package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type FarmReadQueryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadQueryInMemory(s *storage.FarmReadStorage) query.FarmRead {
	return FarmReadQueryInMemory{Storage: s}
}

func (s FarmReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farm := storage.FarmRead{}

		for _, val := range s.Storage.FarmReadMap {
			if val.UID == uid {
				farm = val
			}
		}

		result <- query.Result{Result: farm}

		close(result)
	}()

	return result
}

func (s FarmReadQueryInMemory) FindAll() <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farms := []storage.FarmRead{}
		for _, val := range s.Storage.FarmReadMap {
			farms = append(farms, val)
		}

		result <- query.Result{Result: farms}

		close(result)
	}()

	return result
}

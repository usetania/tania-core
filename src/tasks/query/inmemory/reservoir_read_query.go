package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type ReservoirQueryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirQueryInMemory(s *storage.ReservoirReadStorage) query.ReservoirQuery {
	return ReservoirQueryInMemory{Storage: s}
}

func (s ReservoirQueryInMemory) FindReservoirByID(reservoirUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		ci := query.TaskReservoirQueryResult{}
		for _, val := range s.Storage.ReservoirReadMap {
			// WARNING, domain leakage

			if val.UID == reservoirUID {
				ci.UID = val.UID
				ci.Name = val.Name
			}
		}

		result <- query.QueryResult{Result: ci}

		close(result)
	}()

	return result
}

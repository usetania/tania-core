package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/storage"
	"github.com/usetania/tania-core/src/tasks/query"
)

type ReservoirQueryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirQueryInMemory(s *storage.ReservoirReadStorage) query.Reservoir {
	return ReservoirQueryInMemory{Storage: s}
}

func (s ReservoirQueryInMemory) FindReservoirByID(reservoirUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		ci := query.TaskReservoirResult{}

		for _, val := range s.Storage.ReservoirReadMap {
			// WARNING, domain leakage
			if val.UID == reservoirUID {
				ci.UID = val.UID
				ci.Name = val.Name
			}
		}

		result <- query.Result{Result: ci}

		close(result)
	}()

	return result
}

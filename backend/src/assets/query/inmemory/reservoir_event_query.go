package inmemory

import (
	"sort"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type ReservoirEventQueryInMemory struct {
	Storage *storage.ReservoirEventStorage
}

func NewReservoirEventQueryInMemory(s *storage.ReservoirEventStorage) query.ReservoirEvent {
	return &ReservoirEventQueryInMemory{Storage: s}
}

func (f *ReservoirEventQueryInMemory) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		events := []storage.ReservoirEvent{}

		for _, v := range f.Storage.ReservoirEvents {
			if v.ReservoirUID == uid {
				events = append(events, v)
			}
		}

		sort.Slice(events, func(i, j int) bool {
			return events[i].Version < events[j].Version
		})

		result <- query.Result{Result: events}
	}()

	return result
}

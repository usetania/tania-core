package inmemory

import (
	"sort"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type FarmEventQueryInMemory struct {
	Storage *storage.FarmEventStorage
}

func NewFarmEventQueryInMemory(s *storage.FarmEventStorage) query.FarmEvent {
	return &FarmEventQueryInMemory{Storage: s}
}

func (f *FarmEventQueryInMemory) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		events := []storage.FarmEvent{}

		for _, v := range f.Storage.FarmEvents {
			if v.FarmUID == uid {
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

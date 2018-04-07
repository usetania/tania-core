package inmemory

import (
	"sort"

	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmEventQueryInMemory struct {
	Storage *storage.FarmEventStorage
}

func NewFarmEventQueryInMemory(s *storage.FarmEventStorage) query.FarmEventQuery {
	return &FarmEventQueryInMemory{Storage: s}
}

func (f *FarmEventQueryInMemory) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

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

		result <- query.QueryResult{Result: events}
	}()

	return result
}

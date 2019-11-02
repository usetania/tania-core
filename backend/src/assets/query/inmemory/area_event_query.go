package inmemory

import (
	"sort"

	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type AreaEventQueryInMemory struct {
	Storage *storage.AreaEventStorage
}

func NewAreaEventQueryInMemory(s *storage.AreaEventStorage) query.AreaEventQuery {
	return &AreaEventQueryInMemory{Storage: s}
}

func (f *AreaEventQueryInMemory) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		events := []storage.AreaEvent{}
		for _, v := range f.Storage.AreaEvents {
			if v.AreaUID == uid {
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

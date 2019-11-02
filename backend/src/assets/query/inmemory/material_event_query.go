package inmemory

import (
	"sort"

	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type MaterialEventQueryInMemory struct {
	Storage *storage.MaterialEventStorage
}

func NewMaterialEventQueryInMemory(s *storage.MaterialEventStorage) query.MaterialEventQuery {
	return &MaterialEventQueryInMemory{Storage: s}
}

func (f *MaterialEventQueryInMemory) FindAllByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		events := []storage.MaterialEvent{}
		for _, v := range f.Storage.MaterialEvents {
			if v.MaterialUID == uid {
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

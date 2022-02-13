package inmemory

import (
	"sort"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type MaterialEventQueryInMemory struct {
	Storage *storage.MaterialEventStorage
}

func NewMaterialEventQueryInMemory(s *storage.MaterialEventStorage) query.MaterialEvent {
	return &MaterialEventQueryInMemory{Storage: s}
}

func (f *MaterialEventQueryInMemory) FindAllByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

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

		result <- query.Result{Result: events}
	}()

	return result
}

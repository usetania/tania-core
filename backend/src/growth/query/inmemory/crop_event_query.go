package inmemory

import (
	"sort"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropEventQueryInMemory struct {
	Storage *storage.CropEventStorage
}

func NewCropEventQueryInMemory(s *storage.CropEventStorage) query.CropEventQuery {
	return &CropEventQueryInMemory{Storage: s}
}

func (f *CropEventQueryInMemory) FindAllByCropID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		f.Storage.Lock.RLock()
		defer f.Storage.Lock.RUnlock()

		events := []storage.CropEvent{}

		for _, v := range f.Storage.CropEvents {
			if v.CropUID == uid {
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

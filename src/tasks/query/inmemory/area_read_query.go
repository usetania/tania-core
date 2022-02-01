package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/Tanibox/tania-core/src/tasks/query"
	"github.com/gofrs/uuid"
)

type AreaQueryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaQueryInMemory(s *storage.AreaReadStorage) query.Area {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := query.TaskAreaResult{}

		for _, val := range s.Storage.AreaReadMap {
			if val.UID == uid {
				area.UID = uid
				area.Name = val.Name
			}
		}

		result <- query.Result{Result: area}

		close(result)
	}()

	return result
}

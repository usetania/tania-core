package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type AreaQueryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaQueryInMemory(s *storage.AreaReadStorage) query.AreaQuery {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := query.TaskAreaQueryResult{}
		for _, val := range s.Storage.AreaReadMap {
			if val.UID == uid {
				area.UID = uid
				area.Name = val.Name
			}
		}

		result <- query.QueryResult{Result: area}

		close(result)
	}()

	return result
}

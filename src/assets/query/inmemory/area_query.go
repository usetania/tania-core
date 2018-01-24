package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type AreaQueryInMemory struct {
	Storage *storage.AreaStorage
}

func NewAreaQueryInMemory(s *storage.AreaStorage) query.AreaQuery {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindAreasByReservoirID(reservoirID string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		areas := []domain.Area{}
		for _, val := range s.Storage.AreaMap {
			if val.Reservoir.UID.String() == reservoirID {
				areas = append(areas, val)
			}
		}

		result <- query.QueryResult{Result: areas}

		close(result)
	}()

	return result
}

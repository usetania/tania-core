package query

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type AreaQuery interface {
	FindAreasByReservoirID(reservoirID string) <-chan QueryResult
}

type AreaQueryInMemory struct {
	Storage *storage.AreaStorage
}

func NewAreaQueryInMemory(s *storage.AreaStorage) AreaQuery {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindAreasByReservoirID(reservoirID string) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		areas := []domain.Area{}
		for _, val := range s.Storage.AreaMap {
			if val.Reservoir.UID.String() == reservoirID {
				areas = append(areas, val)
			}
		}

		result <- QueryResult{Result: areas}

		close(result)
	}()

	return result
}

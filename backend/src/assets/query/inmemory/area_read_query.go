package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type AreaReadQueryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaReadQueryInMemory(s *storage.AreaReadStorage) query.AreaReadQuery {
	return AreaReadQueryInMemory{Storage: s}
}

func (s AreaReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := storage.AreaRead{}
		for _, val := range s.Storage.AreaReadMap {
			if val.UID == uid {
				area = val
			}
		}

		result <- query.QueryResult{Result: area}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) FindAllByFarm(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		areas := []storage.AreaRead{}
		for _, val := range s.Storage.AreaReadMap {
			if val.Farm.UID == farmUID {
				areas = append(areas, val)
			}
		}

		result <- query.QueryResult{Result: areas}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) FindByIDAndFarm(areaUID, farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := storage.AreaRead{}
		for _, val := range s.Storage.AreaReadMap {
			if val.Farm.UID == farmUID && val.UID == areaUID {
				area = val
			}
		}

		result <- query.QueryResult{Result: area}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) FindAreasByReservoirID(reservoirUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		areas := []storage.AreaRead{}
		for _, val := range s.Storage.AreaReadMap {
			if val.Reservoir.UID == reservoirUID {
				areas = append(areas, val)
			}
		}

		result <- query.QueryResult{Result: areas}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) CountAreas(farmUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		total := 0
		for _, val := range s.Storage.AreaReadMap {
			if val.Farm.UID == farmUID {
				total++
			}
		}

		result <- query.QueryResult{Result: total}

		close(result)
	}()

	return result
}

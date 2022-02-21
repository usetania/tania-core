package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaReadQueryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaReadQueryInMemory(s *storage.AreaReadStorage) query.AreaRead {
	return AreaReadQueryInMemory{Storage: s}
}

func (s AreaReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := storage.AreaRead{}

		for _, val := range s.Storage.AreaReadMap {
			if val.UID == uid {
				area = val
			}
		}

		result <- query.Result{Result: area}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) FindAllByFarm(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		areas := []storage.AreaRead{}

		for _, val := range s.Storage.AreaReadMap {
			if val.Farm.UID == farmUID {
				areas = append(areas, val)
			}
		}

		result <- query.Result{Result: areas}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) FindByIDAndFarm(areaUID, farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := storage.AreaRead{}

		for _, val := range s.Storage.AreaReadMap {
			if val.Farm.UID == farmUID && val.UID == areaUID {
				area = val
			}
		}

		result <- query.Result{Result: area}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) FindAreasByReservoirID(reservoirUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		areas := []storage.AreaRead{}

		for _, val := range s.Storage.AreaReadMap {
			if val.Reservoir.UID == reservoirUID {
				areas = append(areas, val)
			}
		}

		result <- query.Result{Result: areas}

		close(result)
	}()

	return result
}

func (s AreaReadQueryInMemory) CountAreas(farmUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		total := 0

		for _, val := range s.Storage.AreaReadMap {
			if val.Farm.UID == farmUID {
				total++
			}
		}

		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}

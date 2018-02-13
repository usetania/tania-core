package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type ReservoirReadQueryInMemory struct {
	Storage *storage.ReservoirReadStorage
}

func NewReservoirReadQueryInMemory(s *storage.ReservoirReadStorage) query.ReservoirReadQuery {
	return ReservoirReadQueryInMemory{Storage: s}
}

func (s ReservoirReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		reservoir := query.ReservoirReadQueryResult{}
		for _, val := range s.Storage.ReservoirReadMap {
			if val.UID == uid {
				reservoir.UID = val.UID
				reservoir.Name = val.Name
				reservoir.WaterSource = query.WaterSource{
					Type:     val.WaterSource.Type,
					Capacity: val.WaterSource.Capacity,
				}
				reservoir.FarmUID = val.FarmUID
				reservoir.CreatedDate = val.CreatedDate
			}
		}

		result <- query.QueryResult{Result: reservoir}

		close(result)
	}()

	return result
}

func (s ReservoirReadQueryInMemory) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farms := []query.ReservoirReadQueryResult{}
		for _, val := range s.Storage.ReservoirReadMap {
			farms = append(farms, query.ReservoirReadQueryResult{
				UID:  val.UID,
				Name: val.Name,
				WaterSource: query.WaterSource{
					Type:     val.WaterSource.Type,
					Capacity: val.WaterSource.Capacity,
				},
				FarmUID:     val.FarmUID,
				CreatedDate: val.CreatedDate,
			})
		}

		result <- query.QueryResult{Result: farms}

		close(result)
	}()

	return result
}

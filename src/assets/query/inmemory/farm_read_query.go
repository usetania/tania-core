package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type FarmReadQueryInMemory struct {
	Storage *storage.FarmReadStorage
}

func NewFarmReadQueryInMemory(s *storage.FarmReadStorage) query.FarmReadQuery {
	return FarmReadQueryInMemory{Storage: s}
}

func (s FarmReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farm := storage.FarmRead{}
		for _, val := range s.Storage.FarmReadMap {
			if val.UID == uid {
				farm = val
				// farm.UID = val.UID
				// farm.Name = val.Name
				// farm.Type = val.Type
				// farm.Latitude = val.Latitude
				// farm.Longitude = val.Longitude
				// farm.CountryCode = val.CountryCode
				// farm.CityCode = val.CityCode
				// farm.CreatedDate = val.CreatedDate
			}
		}

		result <- query.QueryResult{Result: farm}

		close(result)
	}()

	return result
}

func (s FarmReadQueryInMemory) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		farms := []storage.FarmRead{}
		for _, val := range s.Storage.FarmReadMap {
			farms = append(farms, val)
			// farms = append(farms, query.FarmReadQueryResult{
			// 	UID:         val.UID,
			// 	Name:        val.Name,
			// 	Type:        val.Type,
			// 	Latitude:    val.Latitude,
			// 	Longitude:   val.Longitude,
			// 	CountryCode: val.CountryCode,
			// 	CityCode:    val.CityCode,
			// 	CreatedDate: val.CreatedDate,
			// })
		}

		result <- query.QueryResult{Result: farms}

		close(result)
	}()

	return result
}

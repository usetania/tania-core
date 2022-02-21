package inmemory

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/storage"
	"github.com/usetania/tania-core/src/growth/query"
)

type AreaReadQueryInMemory struct {
	Storage *storage.AreaReadStorage
}

func NewAreaReadQueryInMemory(s *storage.AreaReadStorage) query.AreaReadQuery {
	return AreaReadQueryInMemory{Storage: s}
}

func (s AreaReadQueryInMemory) FindByID(uid uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := query.CropAreaQueryResult{}

		for _, val := range s.Storage.AreaReadMap {
			if val.UID == uid {
				area.UID = uid
				area.Name = val.Name
				area.Size.Value = val.Size.Value
				area.Size.Symbol = val.Size.Unit.Symbol
				area.Type = val.Type
				area.Location = val.Location.Code
				area.FarmUID = val.Farm.UID
			}
		}

		result <- query.Result{Result: area}

		close(result)
	}()

	return result
}

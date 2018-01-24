package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type AreaQueryInMemory struct {
	Storage *storage.AreaStorage
}

func NewAreaQueryInMemory(s *storage.AreaStorage) query.AreaQuery {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindByID(uid uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := query.CropAreaQueryResult{}
		for _, val := range s.Storage.AreaMap {
			if val.UID == uid {
				area.UID = uid
				area.Name = val.Name
				area.Size.Value = val.Size.Value
				area.Size.Symbol = val.Size.Unit.Symbol
				area.Type = val.Type.Code
				area.Location = val.Location.Code
				area.FarmUID = val.Farm.UID
			}
		}

		result <- query.QueryResult{Result: area}

		close(result)
	}()

	return result
}

package query

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/growth/domain"
	uuid "github.com/satori/go.uuid"
)

type AreaQuery interface {
	FindByID(areaUID uuid.UUID) <-chan QueryResult
}

type AreaQueryInMemory struct {
	Storage *storage.AreaStorage
}

func NewAreaQueryInMemory(s *storage.AreaStorage) AreaQuery {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindByID(uid uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := domain.CropArea{}
		for _, val := range s.Storage.AreaMap {
			if val.UID == uid {
				area.UID = uid
				area.Name = val.Name
				// insert size here
				area.Type = val.Type
				area.Location = val.Location
				area.FarmUID = val.Farm.UID
			}
		}

		result <- QueryResult{Result: area}

		close(result)
	}()

	return result
}

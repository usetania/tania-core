package query

import (
	"github.com/Tanibox/tania-server/src/growth/storage"
	uuid "github.com/satori/go.uuid"
)

type AreaQuery interface {
	FindByID(areaUID uuid.UUID) <-chan QueryResult
}

type AreaQueryInMemory struct {
	Storage *storage.AreaStorage
}

type CropArea struct {
	UID      uuid.UUID
	Name     string       `json:"name"`
	Size     CropAreaUnit `json:"size"`
	Type     string       `json:"type"`
	Location string       `json:"location"`
}

type CropAreaUnit struct {
	Value  float32
	Symbol string
}

func NewAreaQueryInMemory(s *storage.CropStorage) AreaQuery {
	return AreaQueryInMemory{Storage: s}
}

func (s AreaQueryInMemory) FindByID(uid uuid.UUID) <-chan QueryResult {
	result := make(chan QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		area := CropArea{}
		for _, val := range s.Storage.AreaMap {
			if val.UID == uid {
				area.UID = uid
				area.Name = val.Name
				area.Size = CropAreaUnit{
					Value:  val.Size.Value,
					Symbol: val.Symbol(),
				}
				area.Type = val.Type
				area.Location = val.Location
			}
		}

		result <- QueryResult{Result: area}

		close(result)
	}()

	return result
}

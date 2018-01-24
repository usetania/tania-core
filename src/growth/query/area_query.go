package query

import (
	assetsdomain "github.com/Tanibox/tania-server/src/assets/domain"
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

				// WARNING! Domain leakage. Please change this if we have better solution
				cropAreaUnit := domain.CropAreaUnit{}
				switch v := val.Size.(type) {
				case assetsdomain.SquareMeter:
					cropAreaUnit.Value = v.Value
					cropAreaUnit.Symbol = v.Symbol()
				case assetsdomain.Hectare:
					cropAreaUnit.Value = v.Value
					cropAreaUnit.Symbol = v.Symbol()
				}

				area.Size = cropAreaUnit
				area.Type = domain.GetAreaType(val.Type.Code)
				area.Location = val.Location
				area.FarmUID = val.Farm.UID
			}
		}

		result <- QueryResult{Result: area}

		close(result)
	}()

	return result
}

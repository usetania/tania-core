package inmemory

import (
	assetsdomain "github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/Tanibox/tania-core/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type MaterialReadQueryInMemory struct {
	Storage *storage.MaterialReadStorage
}

func NewMaterialReadQueryInMemory(s *storage.MaterialReadStorage) query.MaterialReadQuery {
	return MaterialReadQueryInMemory{Storage: s}
}

func (s MaterialReadQueryInMemory) FindByID(inventoryUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		ci := query.CropMaterialQueryResult{}
		for _, val := range s.Storage.MaterialReadMap {
			if val.UID == inventoryUID {
				ci.UID = val.UID
				ci.Name = val.Name
				ci.TypeCode = val.Type.Code()

				// WARNING, domain leakage
				switch v := val.Type.(type) {
				case assetsdomain.MaterialTypeSeed:
					ci.PlantTypeCode = v.PlantType.Code
				case assetsdomain.MaterialTypePlant:
					ci.PlantTypeCode = v.PlantType.Code
				}
			}
		}

		result <- query.QueryResult{Result: ci}

		close(result)
	}()

	return result
}

func (q MaterialReadQueryInMemory) FindMaterialByPlantTypeCodeAndName(plantTypeCode string, name string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		ci := query.CropMaterialQueryResult{}
		for _, val := range q.Storage.MaterialReadMap {
			// WARNING, domain leakage
			switch v := val.Type.(type) {
			case assetsdomain.MaterialTypeSeed:
				if v.PlantType.Code == plantTypeCode && val.Name == name {
					ci.UID = val.UID
					ci.Name = val.Name
					ci.TypeCode = val.Type.Code()
					ci.PlantTypeCode = v.PlantType.Code
				}
			case assetsdomain.MaterialTypePlant:
				if v.PlantType.Code == plantTypeCode && val.Name == name {
					ci.UID = val.UID
					ci.Name = val.Name
					ci.TypeCode = val.Type.Code()
					ci.PlantTypeCode = v.PlantType.Code
				}
			}
		}

		result <- query.QueryResult{Result: ci}

		close(result)
	}()

	return result
}

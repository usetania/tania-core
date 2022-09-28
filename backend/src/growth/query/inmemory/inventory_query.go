package inmemory

import (
	"github.com/gofrs/uuid"
	assetsdomain "github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/storage"
	"github.com/usetania/tania-core/src/growth/query"
)

type MaterialReadQueryInMemory struct {
	Storage *storage.MaterialReadStorage
}

func NewMaterialReadQueryInMemory(s *storage.MaterialReadStorage) query.MaterialReadQuery {
	return MaterialReadQueryInMemory{Storage: s}
}

func (q MaterialReadQueryInMemory) FindByID(inventoryUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		ci := query.CropMaterialQueryResult{}

		for _, val := range q.Storage.MaterialReadMap {
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

		result <- query.Result{Result: ci}

		close(result)
	}()

	return result
}

func (q MaterialReadQueryInMemory) FindMaterialByPlantTypeCodeAndName(plantTypeCode, name string) <-chan query.Result {
	result := make(chan query.Result)

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

		result <- query.Result{Result: ci}

		close(result)
	}()

	return result
}

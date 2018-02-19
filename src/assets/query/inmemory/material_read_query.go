package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

type MaterialReadQueryInMemory struct {
	Storage *storage.MaterialReadStorage
}

func NewMaterialReadQueryInMemory(s *storage.MaterialReadStorage) query.MaterialReadQuery {
	return &MaterialReadQueryInMemory{Storage: s}
}

func (q *MaterialReadQueryInMemory) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		materials := []storage.MaterialRead{}
		for _, val := range q.Storage.MaterialReadMap {
			materials = append(materials, val)
		}

		result <- query.QueryResult{Result: materials}

		close(result)
	}()

	return result
}

func (q *MaterialReadQueryInMemory) FindByID(materialUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		material := storage.MaterialRead{}
		for _, val := range q.Storage.MaterialReadMap {
			if val.UID == materialUID {
				material = val
			}
		}

		result <- query.QueryResult{Result: material}

		close(result)
	}()

	return result
}

func (q *MaterialReadQueryInMemory) FindAllMaterialByPlantType(plantTypeCode string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		materials := []storage.MaterialRead{}
		for _, val := range q.Storage.MaterialReadMap {
			switch v := val.Type.(type) {
			case domain.MaterialTypeSeed:
				if v.PlantType.Code == plantTypeCode {
					materials = append(materials, val)
				}
			case domain.MaterialTypePlant:
				if v.PlantType.Code == plantTypeCode {
					materials = append(materials, val)
				}
			}
		}

		result <- query.QueryResult{Result: materials}

		close(result)
	}()

	return result
}

func (q *MaterialReadQueryInMemory) FindMaterialByPlantTypeAndName(plantTypeCode string, name string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		material := storage.MaterialRead{}
		for _, val := range q.Storage.MaterialReadMap {
			switch v := val.Type.(type) {
			case domain.MaterialTypeSeed:
				if v.PlantType.Code == plantTypeCode && val.Name == name {
					material = val
				}
			case domain.MaterialTypePlant:
				if v.PlantType.Code == plantTypeCode && val.Name == name {
					material = val
				}
			}
		}

		result <- query.QueryResult{Result: material}

		close(result)
	}()

	return result
}

package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	"github.com/Tanibox/tania-server/src/assets/storage"
)

type MaterialQueryInMemory struct {
	Storage *storage.MaterialStorage
}

func NewMaterialQueryInMemory(s *storage.MaterialStorage) query.MaterialQuery {
	return &MaterialQueryInMemory{Storage: s}
}

func (q *MaterialQueryInMemory) FindAll() <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		materials := []domain.Material{}
		for _, val := range q.Storage.MaterialMap {
			materials = append(materials, val)
		}

		result <- query.QueryResult{Result: materials}

		close(result)
	}()

	return result
}

func (q *MaterialQueryInMemory) FindAllMaterialByPlantType(plantTypeCode string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		materials := []domain.Material{}
		for _, val := range q.Storage.MaterialMap {
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

func (q *MaterialQueryInMemory) FindMaterialByPlantTypeAndName(plantTypeCode string, name string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		material := domain.Material{}
		for _, val := range q.Storage.MaterialMap {
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

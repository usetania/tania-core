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

func (q *MaterialQueryInMemory) FindAllSeedMaterialByPlantType(plantType domain.PlantType) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		materials := []domain.Material{}
		for _, val := range q.Storage.MaterialMap {
			s, ok := val.Type.(domain.MaterialTypeSeed)

			if ok && s.PlantType == plantType {
				materials = append(materials, val)
			}
		}

		result <- query.QueryResult{Result: materials}

		close(result)
	}()

	return result
}

func (q *MaterialQueryInMemory) FindSeedMaterialByPlantTypeAndName(plantType domain.PlantType, name string) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		material := domain.Material{}
		for _, val := range q.Storage.MaterialMap {
			s, ok := val.Type.(domain.MaterialTypeSeed)

			if ok && s.PlantType == plantType && val.Name == name {
				material = val
			}
		}

		result <- query.QueryResult{Result: material}

		close(result)
	}()

	return result
}

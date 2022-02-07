package inmemory

import (
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/gofrs/uuid"
)

type MaterialReadQueryInMemory struct {
	Storage *storage.MaterialReadStorage
}

func NewMaterialReadQueryInMemory(s *storage.MaterialReadStorage) query.MaterialRead {
	return &MaterialReadQueryInMemory{Storage: s}
}

func (q *MaterialReadQueryInMemory) FindAll(materialType, materialTypeDetail string, page, limit int) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		materials := []storage.MaterialRead{}
		for _, val := range q.Storage.MaterialReadMap {
			materials = append(materials, val)
		}

		result <- query.Result{Result: materials}

		close(result)
	}()

	return result
}

func (q MaterialReadQueryInMemory) CountAll(materialType, materialTypeDetail string) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		total := len(q.Storage.MaterialReadMap)

		result <- query.Result{Result: total}

		close(result)
	}()

	return result
}

func (q *MaterialReadQueryInMemory) FindByID(materialUID uuid.UUID) <-chan query.Result {
	result := make(chan query.Result)

	go func() {
		q.Storage.Lock.RLock()
		defer q.Storage.Lock.RUnlock()

		material := storage.MaterialRead{}

		for _, val := range q.Storage.MaterialReadMap {
			if val.UID == materialUID {
				material = val
			}
		}

		result <- query.Result{Result: material}

		close(result)
	}()

	return result
}

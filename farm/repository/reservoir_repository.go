package repository

import (
	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/farm/storage"
	uuid "github.com/satori/go.uuid"
)

// ReservoirRepository is a repository
type ReservoirRepository interface {
	Save(val *entity.Reservoir) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// ReservoirRepositoryInMemory is in-memory AreaRepository db implementation
type ReservoirRepositoryInMemory struct {
	Storage *storage.ReservoirStorage
}

func NewReservoirRepositoryInMemory(s *storage.ReservoirStorage) ReservoirRepository {
	return &ReservoirRepositoryInMemory{Storage: s}
}

// FindAll is to find all
func (r *ReservoirRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		reservoirs := []entity.Reservoir{}
		for _, val := range r.Storage.ReservoirMap {
			reservoirs = append(reservoirs, val)
		}

		result <- RepositoryResult{Result: reservoirs}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (r *ReservoirRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: r.Storage.ReservoirMap[uid]}
	}()

	return result
}

// Save is to save
func (r *ReservoirRepositoryInMemory) Save(val *entity.Reservoir) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.Lock()
		defer r.Storage.Lock.Unlock()

		r.Storage.ReservoirMap[val.UID] = *val

		result <- RepositoryResult{Result: val.UID}

		close(result)
	}()

	return result
}

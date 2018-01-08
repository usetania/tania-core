package repository

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
	uuid "github.com/satori/go.uuid"
)

// ReservoirRepository is a repository
type ReservoirRepository interface {
	Save(val *entity.Reservoir) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// ReservoirRepositoryInMemory is in-memory ReservoirRepository db implementation
type ReservoirRepositoryInMemory struct {
	lock         sync.RWMutex
	ReservoirMap map[uuid.UUID]entity.Reservoir
}

func NewReservoirRepositoryInMemory() ReservoirRepository {
	return &ReservoirRepositoryInMemory{ReservoirMap: make(map[uuid.UUID]entity.Reservoir)}
}

// FindAll is to find all
func (r *ReservoirRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.lock.RLock()
		defer r.lock.RUnlock()

		reservoirs := []entity.Reservoir{}
		for _, val := range r.ReservoirMap {
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
		r.lock.RLock()
		defer r.lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: r.ReservoirMap[uid]}
	}()

	return result
}

// Save is to save
func (r *ReservoirRepositoryInMemory) Save(val *entity.Reservoir) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.ReservoirMap[val.UID] = *val

		result <- RepositoryResult{Result: val.UID}

		close(result)
	}()

	return result
}

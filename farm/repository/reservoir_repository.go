package repository

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
)

// ReservoirRepository is a repository
type ReservoirRepository interface {
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
	Count() <-chan RepositoryResult
	Save(val *entity.Reservoir) <-chan RepositoryResult
}

// ReservoirRepositoryInMemory is in-memory ReservoirRepository db implementation
type ReservoirRepositoryInMemory struct {
	lock         sync.RWMutex
	ReservoirMap map[string]entity.Reservoir
}

func NewReservoirRepositoryInMemory() ReservoirRepository {
	return &ReservoirRepositoryInMemory{ReservoirMap: make(map[string]entity.Reservoir)}
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

		uid := getRandomUID()
		val.UID = uid

		r.ReservoirMap[uid] = *val

		result <- RepositoryResult{Result: uid}

		close(result)
	}()

	return result
}

// Count is to count
func (r *ReservoirRepositoryInMemory) Count() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.lock.RLock()
		defer r.lock.RUnlock()

		count := len(r.ReservoirMap)

		result <- RepositoryResult{Result: count}
	}()

	return result
}

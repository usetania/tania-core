package repository

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
	uuid "github.com/satori/go.uuid"
)

// AreaRepository is a repository
type AreaRepository interface {
	Save(val *entity.Area) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// AreaRepositoryInMemory is in-memory AreaRepository db implementation
type AreaRepositoryInMemory struct {
	lock    sync.RWMutex
	AreaMap map[uuid.UUID]entity.Area
}

func NewAreaRepositoryInMemory() AreaRepository {
	return &AreaRepositoryInMemory{AreaMap: make(map[uuid.UUID]entity.Area)}
}

// FindAll is to find all
func (r *AreaRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.lock.RLock()
		defer r.lock.RUnlock()

		areas := []entity.Area{}
		for _, val := range r.AreaMap {
			areas = append(areas, val)
		}

		result <- RepositoryResult{Result: areas}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (r *AreaRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.lock.RLock()
		defer r.lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: r.AreaMap[uid]}
	}()

	return result
}

// Save is to save
func (r *AreaRepositoryInMemory) Save(val *entity.Area) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.AreaMap[val.UID] = *val

		result <- RepositoryResult{Result: val.UID}

		close(result)
	}()

	return result
}

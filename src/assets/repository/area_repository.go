package repository

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

// AreaRepository is a repository
type AreaRepository interface {
	Save(val *domain.Area) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// AreaRepositoryInMemory is in-memory AreaRepository db implementation
type AreaRepositoryInMemory struct {
	Storage *storage.AreaStorage
}

func NewAreaRepositoryInMemory(s *storage.AreaStorage) AreaRepository {
	return &AreaRepositoryInMemory{Storage: s}
}

// FindAll is to find all
func (r *AreaRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		areas := []domain.Area{}
		for _, val := range r.Storage.AreaMap {
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
		r.Storage.Lock.RLock()
		defer r.Storage.Lock.RUnlock()

		uid, err := uuid.FromString(uid)
		if err != nil {
			result <- RepositoryResult{Error: err}
		}

		result <- RepositoryResult{Result: r.Storage.AreaMap[uid]}
	}()

	return result
}

// Save is to save
func (r *AreaRepositoryInMemory) Save(val *domain.Area) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.Storage.Lock.Lock()
		defer r.Storage.Lock.Unlock()

		r.Storage.AreaMap[val.UID] = *val

		result <- RepositoryResult{Error: nil}

		close(result)
	}()

	return result
}

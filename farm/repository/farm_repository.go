package repository

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
)

type FarmRepository interface {
	Save(val *entity.Farm) <-chan RepositoryResult
	FindAll() <-chan RepositoryResult
	FindByID(uid string) <-chan RepositoryResult
}

// ReservoirRepositoryInMemory is in-memory ReservoirRepository db implementation
type FarmRepositoryInMemory struct {
	lock    sync.RWMutex
	FarmMap map[string]entity.Farm
}

func NewFarmRepositoryInMemory() FarmRepository {
	return &FarmRepositoryInMemory{FarmMap: make(map[string]entity.Farm)}
}

func (f *FarmRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.lock.RLock()
		defer f.lock.RUnlock()

		farms := []entity.Farm{}
		for _, val := range f.FarmMap {
			farms = append(farms, val)
		}

		result <- RepositoryResult{Result: farms}

		close(result)
	}()

	return result
}

// Save is to save
func (f *FarmRepositoryInMemory) Save(val *entity.Farm) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.FarmMap[val.UID] = *val

		result <- RepositoryResult{Result: val.UID}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (f *FarmRepositoryInMemory) FindByID(uid string) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.lock.RLock()
		defer f.lock.RUnlock()

		result <- RepositoryResult{Result: f.FarmMap[uid]}
	}()

	return result
}

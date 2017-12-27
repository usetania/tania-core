package repository

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
)

type FarmRepository interface {
	Count() <-chan RepositoryResult
	Save(val *entity.Farm) <-chan RepositoryResult
	Update(val *entity.Farm) <-chan RepositoryResult
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

// Save is to save
func (f *FarmRepositoryInMemory) Save(val *entity.Farm) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.lock.Lock()
		defer f.lock.Unlock()

		uid := getRandomUID()
		val.UID = uid

		f.FarmMap[uid] = *val

		result <- RepositoryResult{Result: uid}

		close(result)
	}()

	return result
}

// Update is to update
func (f *FarmRepositoryInMemory) Update(val *entity.Farm) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.FarmMap[val.UID] = *val

		result <- RepositoryResult{Result: val}

		close(result)
	}()

	return result
}

// Count is to count
func (f *FarmRepositoryInMemory) Count() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		f.lock.RLock()
		defer f.lock.RUnlock()

		count := len(f.FarmMap)

		result <- RepositoryResult{Result: count}
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

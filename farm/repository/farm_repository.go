package repository

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/Tanibox/tania-server/farm"
)

type FarmRepository interface {
	Count() <-chan RepositoryResult
	Save(val *farm.Farm) <-chan RepositoryResult
}

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

// ReservoirRepositoryInMemory is in-memory ReservoirRepository db implementation
type FarmRepositoryInMemory struct {
	lock    sync.RWMutex
	FarmMap map[string]farm.Farm
}

func NewFarmRepositoryInMemory() FarmRepository {
	return &FarmRepositoryInMemory{FarmMap: make(map[string]farm.Farm)}
}

// Save is to save
func (f *FarmRepositoryInMemory) Save(val *farm.Farm) <-chan RepositoryResult {
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

func getRandomUID() string {
	rand.Seed(time.Now().UnixNano())
	uid := rand.Int()
	return strconv.Itoa(uid)
}

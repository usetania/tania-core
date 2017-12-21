package repository

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/Tanibox/tania-server/reservoir"
)

// ReservoirRepository is a repository
type ReservoirRepository interface {
	FindAll() <-chan RepositoryResult
	FindByID(uid string) (reservoir.Reservoir, error)
	Save(val *reservoir.Reservoir) <-chan RepositoryResult
}

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

// ReservoirRepositoryInMemory is in-memory ReservoirRepository db implementation
type ReservoirRepositoryInMemory struct {
	mutex        *sync.Mutex
	ReservoirMap map[string]reservoir.Reservoir
}

// FindAll is to find all
func (r *ReservoirRepositoryInMemory) FindAll() <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		reservoirs := []reservoir.Reservoir{}
		for _, val := range r.ReservoirMap {
			reservoirs = append(reservoirs, val)
		}

		result <- RepositoryResult{Result: reservoirs}

		close(result)
	}()

	return result
}

// FindByID is to find by ID
func (r *ReservoirRepositoryInMemory) FindByID(uid string) (reservoir.Reservoir, error) {
	// Unimplemented
	return r.ReservoirMap[uid], nil
}

// Save is to save
func (r *ReservoirRepositoryInMemory) Save(val *reservoir.Reservoir) <-chan RepositoryResult {
	result := make(chan RepositoryResult)

	go func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		if r.ReservoirMap == nil {
			r.ReservoirMap = map[string]reservoir.Reservoir{}
		}

		uid := getRandomUID()

		val.UID = uid
		r.ReservoirMap[uid] = *val

		result <- RepositoryResult{Result: uid}

		close(result)
	}()

	return result
}

func getRandomUID() string {
	rand.Seed(time.Now().UnixNano())
	uid := rand.Int()
	return strconv.Itoa(uid)
}

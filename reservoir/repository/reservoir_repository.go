package repository

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/Tanibox/tania-server/reservoir"
)

// ReservoirRepository is a repository
type ReservoirRepository interface {
	FindAll() ([]reservoir.Reservoir, error)
	FindByID(uid string) (reservoir.Reservoir, error)
	Save(val reservoir.Reservoir) (string, error)
}

// ReservoirRepositoryInMemory is in-memory ReservoirRepository db implementation
type ReservoirRepositoryInMemory struct {
	ReservoirMap map[string]reservoir.Reservoir
}

// FindAll is to find all
func (r ReservoirRepositoryInMemory) FindAll() ([]reservoir.Reservoir, error) {
	reservoirs := []reservoir.Reservoir{}

	for _, val := range r.ReservoirMap {
		reservoirs = append(reservoirs, val)
	}

	return reservoirs, nil
}

// FindByID is to find by ID
func (r ReservoirRepositoryInMemory) FindByID(uid string) (reservoir.Reservoir, error) {
	return r.ReservoirMap[uid], nil
}

// Save is to save
func (r *ReservoirRepositoryInMemory) Save(val reservoir.Reservoir) (UID string, err error) {
	if r.ReservoirMap == nil {
		r.ReservoirMap = map[string]reservoir.Reservoir{}
	}

	uid := getRandomUID()

	val.UID = uid
	r.ReservoirMap[uid] = val

	return uid, nil
}

func getRandomUID() string {
	rand.Seed(time.Now().UnixNano())
	uid := rand.Int()
	return strconv.Itoa(uid)
}

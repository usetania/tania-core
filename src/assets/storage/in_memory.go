package storage

import (
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type FarmStorage struct {
	Lock    *deadlock.RWMutex
	FarmMap map[uuid.UUID]domain.Farm
}

func CreateFarmStorage() *FarmStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("DEADLOCK!")
	}

	return &FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm), Lock: &rwMutex}
}

type AreaStorage struct {
	Lock    *deadlock.RWMutex
	AreaMap map[uuid.UUID]domain.Area
}

func CreateAreaStorage() *AreaStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("DEADLOCK!")
	}

	return &AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area), Lock: &rwMutex}
}

type ReservoirStorage struct {
	Lock         *deadlock.RWMutex
	ReservoirMap map[uuid.UUID]domain.Reservoir
}

func CreateReservoirStorage() *ReservoirStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("DEADLOCK!")
	}

	return &ReservoirStorage{ReservoirMap: make(map[uuid.UUID]domain.Reservoir), Lock: &rwMutex}
}

type InventoryMaterialStorage struct {
	Lock                 *deadlock.RWMutex
	InventoryMaterialMap map[uuid.UUID]domain.InventoryMaterial
}

func CreateInventoryMaterialStorage() *InventoryMaterialStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("DEADLOCK!")
	}

	return &InventoryMaterialStorage{InventoryMaterialMap: make(map[uuid.UUID]domain.InventoryMaterial), Lock: &rwMutex}
}

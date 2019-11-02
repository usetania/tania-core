package storage

import (
	"fmt"
	"time"

	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type FarmEventStorage struct {
	Lock       *deadlock.RWMutex
	FarmEvents []FarmEvent
}

func CreateFarmEventStorage() *FarmEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("FARM EVENT STORAGE DEADLOCK!")
	}

	return &FarmEventStorage{Lock: &rwMutex}
}

type FarmReadStorage struct {
	Lock        *deadlock.RWMutex
	FarmReadMap map[uuid.UUID]FarmRead
}

func CreateFarmReadStorage() *FarmReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("FARM READ STORAGE DEADLOCK!")
	}

	return &FarmReadStorage{FarmReadMap: make(map[uuid.UUID]FarmRead), Lock: &rwMutex}
}

type ReservoirEventStorage struct {
	Lock            *deadlock.RWMutex
	ReservoirEvents []ReservoirEvent
}

func CreateReservoirEventStorage() *ReservoirEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("RESERVOIR EVENT STORAGE DEADLOCK!")
	}

	return &ReservoirEventStorage{Lock: &rwMutex}
}

type ReservoirReadStorage struct {
	Lock             *deadlock.RWMutex
	ReservoirReadMap map[uuid.UUID]ReservoirRead
}

func CreateReservoirReadStorage() *ReservoirReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("RESERVOIR READ STORAGE DEADLOCK!")
	}

	return &ReservoirReadStorage{ReservoirReadMap: make(map[uuid.UUID]ReservoirRead), Lock: &rwMutex}
}

type AreaEventStorage struct {
	Lock       *deadlock.RWMutex
	AreaEvents []AreaEvent
}

func CreateAreaEventStorage() *AreaEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("AREA EVENT STORAGE DEADLOCK!")
	}

	return &AreaEventStorage{Lock: &rwMutex}
}

type AreaReadStorage struct {
	Lock        *deadlock.RWMutex
	AreaReadMap map[uuid.UUID]AreaRead
}

func CreateAreaReadStorage() *AreaReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("Area READ STORAGE DEADLOCK!")
	}

	return &AreaReadStorage{AreaReadMap: make(map[uuid.UUID]AreaRead), Lock: &rwMutex}
}

type MaterialEventStorage struct {
	Lock           *deadlock.RWMutex
	MaterialEvents []MaterialEvent
}

func CreateMaterialEventStorage() *MaterialEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("MATERIAL EVENT STORAGE DEADLOCK!")
	}

	return &MaterialEventStorage{Lock: &rwMutex}
}

type MaterialReadStorage struct {
	Lock            *deadlock.RWMutex
	MaterialReadMap map[uuid.UUID]MaterialRead
}

func CreateMaterialReadStorage() *MaterialReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("MATERIAL READ STORAGE DEADLOCK!")
	}

	return &MaterialReadStorage{MaterialReadMap: make(map[uuid.UUID]MaterialRead), Lock: &rwMutex}
}

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
		fmt.Println("FARM STORAGE DEADLOCK!")
	}

	return &FarmStorage{FarmMap: make(map[uuid.UUID]domain.Farm), Lock: &rwMutex}
}

type FarmEventStorage struct {
	Lock       *deadlock.RWMutex
	FarmEvents []FarmEvent
}

type FarmEvent struct {
	FarmUID uuid.UUID
	Version int
	Event   interface{}
}

func CreateFarmEventStorage() *FarmEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("FARM EVENT STORAGE DEADLOCK!")
	}

	return &FarmEventStorage{Lock: &rwMutex}
}

type FarmRead struct {
	UID         uuid.UUID `json:"uid"`
	Name        string    `json:"name"`
	Latitude    string    `json:"latitude"`
	Longitude   string    `json:"longitude"`
	Type        string    `json:"type"`
	CountryCode string    `json:"country_code"`
	CityCode    string    `json:"city_code"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
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

type AreaStorage struct {
	Lock    *deadlock.RWMutex
	AreaMap map[uuid.UUID]domain.Area
}

func CreateAreaStorage() *AreaStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("AREA STORAGE DEADLOCK!")
	}

	return &AreaStorage{AreaMap: make(map[uuid.UUID]domain.Area), Lock: &rwMutex}
}

type AreaEventStorage struct {
	Lock       *deadlock.RWMutex
	AreaEvents []AreaEvent
}

type AreaEvent struct {
	AreaUID uuid.UUID
	Version int
	Event   interface{}
}

func CreateAreaEventStorage() *AreaEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("AREA EVENT STORAGE DEADLOCK!")
	}

	return &AreaEventStorage{Lock: &rwMutex}
}

type AreaRead struct {
	UID         uuid.UUID     `json:"uid"`
	Name        string        `json:"name"`
	Size        AreaSize      `json:"size"`
	Location    AreaLocation  `json:"location"`
	Type        AreaType      `json:"type"`
	Photo       AreaPhoto     `json:"photo"`
	CreatedDate time.Time     `json:"created_date"`
	Notes       []AreaNote    `json:"notes"`
	Farm        AreaFarm      `json:"farm"`
	Reservoir   AreaReservoir `json:"reservoir"`
}

type AreaFarm struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type AreaReservoir struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type AreaSize domain.AreaSize
type AreaLocation domain.AreaLocation
type AreaType domain.AreaType
type AreaPhoto domain.AreaPhoto
type AreaNote domain.AreaNote

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

type ReservoirStorage struct {
	Lock         *deadlock.RWMutex
	ReservoirMap map[uuid.UUID]domain.Reservoir
}

func CreateReservoirStorage() *ReservoirStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("RESERVOIR STORAGE DEADLOCK!")
	}

	return &ReservoirStorage{ReservoirMap: make(map[uuid.UUID]domain.Reservoir), Lock: &rwMutex}
}

type ReservoirEventStorage struct {
	Lock            *deadlock.RWMutex
	ReservoirEvents []ReservoirEvent
}

type ReservoirEvent struct {
	ReservoirUID uuid.UUID
	Version      int
	Event        interface{}
}

func CreateReservoirEventStorage() *ReservoirEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("RESERVOIR EVENT STORAGE DEADLOCK!")
	}

	return &ReservoirEventStorage{Lock: &rwMutex}
}

type ReservoirRead struct {
	UID             uuid.UUID       `json:"uid"`
	Name            string          `json:"name"`
	WaterSource     WaterSource     `json:"water_source"`
	Farm            ReservoirFarm   `json:"farm"`
	Notes           []ReservoirNote `json:"notes"`
	CreatedDate     time.Time       `json:"created_date"`
	InstalledToArea []AreaInstalled `json:"installed_to_area"`
}

type WaterSource struct {
	Type     string
	Capacity float32
}

type ReservoirFarm struct {
	UID  uuid.UUID
	Name string
}

type ReservoirNote domain.ReservoirNote

type AreaInstalled struct {
	UID  uuid.UUID
	Name string
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

type MaterialStorage struct {
	Lock        *deadlock.RWMutex
	MaterialMap map[uuid.UUID]domain.Material
}

func CreateMaterialStorage() *MaterialStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("MATERIAL STORAGE DEADLOCK!")
	}

	return &MaterialStorage{MaterialMap: make(map[uuid.UUID]domain.Material), Lock: &rwMutex}
}

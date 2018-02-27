package storage

import (
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type FarmEventStorage struct {
	Lock       *deadlock.RWMutex
	FarmEvents []FarmEvent
}

type FarmEvent struct {
	FarmUID     uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
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

type ReservoirEventStorage struct {
	Lock            *deadlock.RWMutex
	ReservoirEvents []ReservoirEvent
}

type ReservoirEvent struct {
	ReservoirUID uuid.UUID
	Version      int
	CreatedDate  time.Time
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
	Type     string  `json:"type"`
	Capacity float32 `json:"capacity"`
}

type ReservoirFarm struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
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

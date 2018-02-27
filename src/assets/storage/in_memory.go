package storage

import (
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

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
	AreaUID     uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
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
	Type        string        `json:"type"`
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

type MaterialEventStorage struct {
	Lock           *deadlock.RWMutex
	MaterialEvents []MaterialEvent
}

type MaterialEvent struct {
	MaterialUID uuid.UUID
	Version     int
	Event       interface{}
}

func CreateMaterialEventStorage() *MaterialEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("MATERIAL EVENT STORAGE DEADLOCK!")
	}

	return &MaterialEventStorage{Lock: &rwMutex}
}

type MaterialRead struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	PricePerUnit   Money            `json:"price_per_unit"`
	Type           MaterialType     `json:"type"`
	Quantity       MaterialQuantity `json:"quantity"`
	ExpirationDate *time.Time       `json:"expiration_date"`
	Notes          *string          `json:"notes"`
	IsExpense      *bool            `json:"is_expense"`
	ProducedBy     *string          `json:"produced_by"`
	CreatedDate    time.Time        `json:"created_date"`
}

type Money domain.Money
type MaterialType domain.MaterialType
type MaterialQuantity domain.MaterialQuantity

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

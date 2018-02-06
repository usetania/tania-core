package storage

import (
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type CropStorage struct {
	Lock    *deadlock.RWMutex
	CropMap map[uuid.UUID]domain.Crop
}

func CreateCropStorage() *CropStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP STORAGE DEADLOCK!")
	}

	return &CropStorage{CropMap: make(map[uuid.UUID]domain.Crop), Lock: &rwMutex}
}

type CropEventStorage struct {
	Lock         *deadlock.RWMutex
	CropEventMap map[uuid.UUID]interface{}
}

func CreateCropEventStorage() *CropEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP EVENT STORAGE DEADLOCK!")
	}

	return &CropEventStorage{CropEventMap: make(map[uuid.UUID]interface{}), Lock: &rwMutex}
}

type CropList struct {
	UID          uuid.UUID
	VarietyName  string
	BatchID      string
	InventoryUID uuid.UUID
	InitialArea  InitialArea
	MovedArea    []MovedArea
	AreaStatus   AreaStatus
	CreatedDate  time.Time
	FarmUID      uuid.UUID
}

type InitialArea struct {
	AreaUID         uuid.UUID
	Name            string
	InitialQuantity Container
	CurrentQuantity Container
	LastWatered     *time.Time
}

type MovedArea struct {
	AreaUID         uuid.UUID
	Name            string
	CurrentQuantity Container
	LastWatered     *time.Time
}

type Container struct {
	Type     string
	Quantity int
	Cell     int
}

type AreaStatus struct {
	Seeding int
	Growing int
	Dumped  int
}

type CropListStorage struct {
	Lock        *deadlock.RWMutex
	CropListMap map[uuid.UUID]CropList
}

func CreateCropListStorage() *CropListStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP LIST STORAGE DEADLOCK!")
	}

	return &CropListStorage{CropListMap: make(map[uuid.UUID]CropList), Lock: &rwMutex}
}

const (
	SeedActivityCode = "SEED"
)

type CropActivity struct {
	UID          uuid.UUID
	ActivityType ActivityType
	CreatedDate  time.Time
	Description  string
}

type ActivityType interface {
	Code() string
}

type SeedActivity struct {
	Quantity      int
	ContainerType string
	AreaUID       uuid.UUID
	AreaName      string
	BatchID       string
}

func (a SeedActivity) Code() string {
	return SeedActivityCode
}

type CropActivityStorage struct {
	Lock            *deadlock.RWMutex
	CropActivityMap map[uuid.UUID]CropActivity
}

func CreateCropActivityStorage() *CropActivityStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP LIST STORAGE DEADLOCK!")
	}

	return &CropActivityStorage{CropActivityMap: make(map[uuid.UUID]CropActivity), Lock: &rwMutex}
}

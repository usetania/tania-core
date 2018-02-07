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
	BatchID      string
	Status       string
	Type         string
	Container    Container
	InventoryUID uuid.UUID
	VarietyName  string
	FarmUID      uuid.UUID
	FarmName     string
	CreatedDate  time.Time
	Photos       []domain.CropPhoto
	AreaStatus   AreaStatus

	// Fields to track crop's movement
	InitialArea      InitialArea
	MovedArea        []MovedArea
	HarvestedStorage []HarvestedStorage
	Trash            []Trash

	// Notes
	Notes map[uuid.UUID]domain.CropNote
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

type HarvestedStorage struct {
	Quantity             int
	ProducedGramQuantity float32
	SourceAreaUID        uuid.UUID
	SourceAreaName       string
	CreatedDate          time.Time
	LastUpdated          time.Time
}

type Trash struct {
	Quantity       int
	SourceAreaUID  uuid.UUID
	SourceAreaName string
	CreatedDate    time.Time
	LastUpdated    time.Time
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
	SeedActivityCode    = "SEED"
	MoveActivityCode    = "MOVE"
	HarvestActivityCode = "HARVEST"
	DumpActivityCode    = "DUMP"
)

type CropActivity struct {
	UID           uuid.UUID
	BatchID       string
	VarietyName   string
	ContainerType string
	ActivityType  ActivityType
	CreatedDate   time.Time
	Description   string
}

type ActivityType interface {
	Code() string
}

type SeedActivity struct {
	AreaUID  uuid.UUID
	AreaName string
	Quantity int
}

func (a SeedActivity) Code() string {
	return SeedActivityCode
}

type MoveActivity struct {
	SrcAreaUID  uuid.UUID
	SrcAreaName string
	DstAreaUID  uuid.UUID
	DstAreaName string
	Quantity    int
}

func (a MoveActivity) Code() string {
	return MoveActivityCode
}

type CropActivityStorage struct {
	Lock            *deadlock.RWMutex
	CropActivityMap []CropActivity
}

func CreateCropActivityStorage() *CropActivityStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP LIST STORAGE DEADLOCK!")
	}

	return &CropActivityStorage{CropActivityMap: []CropActivity{}, Lock: &rwMutex}
}

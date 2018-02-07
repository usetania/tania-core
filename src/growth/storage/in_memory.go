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
	Lock       *deadlock.RWMutex
	CropEvents []CropEvent
}

type CropEvent struct {
	CropUID uuid.UUID
	Version int
	Event   interface{}
}

func CreateCropEventStorage() *CropEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP EVENT STORAGE DEADLOCK!")
	}

	return &CropEventStorage{Lock: &rwMutex}
}

type CropRead struct {
	UID        uuid.UUID  `json:"uid"`
	BatchID    string     `json:"batch_id"`
	Status     string     `json:"status"`
	Type       string     `json:"type"`
	Container  Container  `json:"container"`
	Inventory  Inventory  `json:"inventory"`
	AreaStatus AreaStatus `json:"area_status"`
	Photos     []domain.CropPhoto
	FarmUID    uuid.UUID

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
	InitialQuantity int
	CurrentQuantity int
	LastWatered     *time.Time `json:"last_watered"`
	LastFertilized  *time.Time `json:"last_fertilized"`
	LastPruned      *time.Time `json:"last_pruned"`
	LastPesticided  *time.Time `json:"last_pesticided"`
}

type MovedArea struct {
	AreaUID         uuid.UUID
	Name            string
	CurrentQuantity int
	LastWatered     *time.Time
	LastFertilized  *time.Time `json:"last_fertilized"`
	LastPruned      *time.Time `json:"last_pruned"`
	LastPesticided  *time.Time `json:"last_pesticided"`
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

type Inventory struct {
	UID       uuid.UUID `json:"uid"`
	PlantType string    `json:"plant_type"`
	Name      string    `json:"name"`
}

type CropReadStorage struct {
	Lock        *deadlock.RWMutex
	CropReadMap map[uuid.UUID]CropRead
}

func CreateCropReadStorage() *CropReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP READ STORAGE DEADLOCK!")
	}

	return &CropReadStorage{CropReadMap: make(map[uuid.UUID]CropRead), Lock: &rwMutex}
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

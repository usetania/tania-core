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
	UID        uuid.UUID   `json:"uid"`
	BatchID    string      `json:"batch_id"`
	Status     string      `json:"status"`
	Type       string      `json:"type"`
	Container  Container   `json:"container"`
	Inventory  Inventory   `json:"inventory"`
	AreaStatus AreaStatus  `json:"area_status"`
	Photos     []CropPhoto `json:"photos"`
	FarmUID    uuid.UUID   `json:"farm_id"`

	// Fields to track crop's movement
	InitialArea      InitialArea        `json:"initial_area"`
	MovedArea        []MovedArea        `json:"moved_area"`
	HarvestedStorage []HarvestedStorage `json:"harvested_storage"`
	Trash            []Trash            `json:"trash"`

	// Notes
	Notes []domain.CropNote `json:"notes"`
}

type InitialArea struct {
	AreaUID         uuid.UUID  `json:"area_id"`
	Name            string     `json:"name"`
	InitialQuantity int        `json:"initial_quantity"`
	CurrentQuantity int        `json:"current_quantity"`
	LastWatered     *time.Time `json:"last_watered"`
	LastFertilized  *time.Time `json:"last_fertilized"`
	LastPruned      *time.Time `json:"last_pruned"`
	LastPesticided  *time.Time `json:"last_pesticided"`
	CreatedDate     time.Time  `json:"created_date"`
	LastUpdated     time.Time  `json:"last_updated"`
}

type MovedArea struct {
	AreaUID         uuid.UUID  `json:"area_id"`
	Name            string     `json:"name"`
	InitialQuantity int        `json:"initial_quantity"`
	CurrentQuantity int        `json:"current_quantity"`
	LastWatered     *time.Time `json:"last_watered"`
	LastFertilized  *time.Time `json:"last_fertilized"`
	LastPruned      *time.Time `json:"last_pruned"`
	LastPesticided  *time.Time `json:"last_pesticided"`
	CreatedDate     time.Time  `json:"created_date"`
	LastUpdated     time.Time  `json:"last_updated"`
}

type HarvestedStorage struct {
	Quantity             int       `json:"quantity"`
	ProducedGramQuantity float32   `json:"produced_gram_quantity"`
	SourceAreaUID        uuid.UUID `json:"source_area_id"`
	SourceAreaName       string    `json:"source_area_name"`
	CreatedDate          time.Time `json:"created_date"`
	LastUpdated          time.Time `json:"last_updated"`
}

type Trash struct {
	Quantity       int       `json:"quantity"`
	SourceAreaUID  uuid.UUID `json:"source_area_id"`
	SourceAreaName string    `json:"source_area_name"`
	CreatedDate    time.Time `json:"created_date"`
	LastUpdated    time.Time `json:"last_updated"`
}

type Container struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Cell     int    `json:"cell"`
}

type AreaStatus struct {
	Seeding int `json:"seeding"`
	Growing int `json:"growing"`
	Dumped  int `json:"dumped"`
}

type Inventory struct {
	UID       uuid.UUID `json:"uid"`
	PlantType string    `json:"plant_type"`
	Name      string    `json:"name"`
}

type CropPhoto struct {
	UID         uuid.UUID `json:"uid"`
	Filename    string    `json:"filename"`
	MimeType    string    `json:"mime_type"`
	Size        int       `json:"size"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Description string    `json:"description"`
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
	PhotoActivityCode   = "PHOTO"
	WaterActivityCode   = "WATER"
)

type CropActivity struct {
	UID           uuid.UUID    `json:"uid"`
	BatchID       string       `json:"batch_id"`
	ContainerType string       `json:"container_type"`
	ActivityType  ActivityType `json:"activity_type"`
	CreatedDate   time.Time    `json:"created_date"`
	Description   string       `json:"description"`
}

type ActivityType interface {
	Code() string
}

type SeedActivity struct {
	AreaUID     uuid.UUID `json:"area_id"`
	AreaName    string    `json:"area_name"`
	Quantity    int       `json:"quantity"`
	SeedingDate time.Time `json:"seeding_date"`
}

func (a SeedActivity) Code() string {
	return SeedActivityCode
}

type MoveActivity struct {
	SrcAreaUID  uuid.UUID `json:"source_area_id"`
	SrcAreaName string    `json:"source_area_name"`
	DstAreaUID  uuid.UUID `json:"destination_area_id"`
	DstAreaName string    `json:"destination_area_name"`
	Quantity    int       `json:"quantity"`
	MovedDate   time.Time `json:"moved_date"`
}

func (a MoveActivity) Code() string {
	return MoveActivityCode
}

type HarvestActivity struct {
	SrcAreaUID           uuid.UUID `json:"source_area_id"`
	SrcAreaName          string    `json:"source_area_name"`
	Quantity             int       `json:"quantity"`
	ProducedGramQuantity float32   `json:"produced_gram_quantity"`
	HarvestDate          time.Time `json:"harvest_date"`
}

func (a HarvestActivity) Code() string {
	return HarvestActivityCode
}

type DumpActivity struct {
	SrcAreaUID  uuid.UUID `json:"source_area_id"`
	SrcAreaName string    `json:"source_area_name"`
	Quantity    int       `json:"quantity"`
	DumpDate    time.Time `json:"dump_date"`
}

func (a DumpActivity) Code() string {
	return DumpActivityCode
}

type WaterActivity struct {
	AreaUID      uuid.UUID `json:"area_id"`
	AreaName     string    `json:"area_name"`
	WateringDate time.Time `json:"watering_date"`
}

func (a WaterActivity) Code() string {
	return WaterActivityCode
}

type PhotoActivity struct {
	UID         uuid.UUID `json:"uid"`
	Filename    string    `json:"filename"`
	MimeType    string    `json:"mime_type"`
	Size        int       `json:"size"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Description string    `json:"description"`
}

func (a PhotoActivity) Code() string {
	return PhotoActivityCode
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

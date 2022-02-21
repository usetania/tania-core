package storage

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
)

type FarmEvent struct {
	FarmUID     uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
}

type FarmRead struct {
	UID         uuid.UUID `json:"uid"`
	Name        string    `json:"name"`
	Latitude    string    `json:"latitude"`
	Longitude   string    `json:"longitude"`
	Type        string    `json:"type"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	IsActive    bool      `json:"is_active"`
	CreatedDate time.Time `json:"created_date"`
}

type ReservoirEvent struct {
	ReservoirUID uuid.UUID
	Version      int
	CreatedDate  time.Time
	Event        interface{}
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
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type AreaEvent struct {
	AreaUID     uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
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

type (
	AreaSize     domain.AreaSize
	AreaLocation domain.AreaLocation
	AreaType     domain.AreaType
	AreaPhoto    domain.AreaPhoto
	AreaNote     domain.AreaNote
)

type MaterialEvent struct {
	MaterialUID uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
}

type MaterialRead struct {
	UID            uuid.UUID        `json:"uid"`
	Name           string           `json:"name"`
	PricePerUnit   PricePerUnit     `json:"price_per_unit"`
	Type           MaterialType     `json:"type"`
	Quantity       MaterialQuantity `json:"quantity"`
	ExpirationDate *time.Time       `json:"expiration_date"`
	Notes          *string          `json:"notes"`
	IsExpense      *bool            `json:"is_expense"`
	ProducedBy     *string          `json:"produced_by"`
	CreatedDate    time.Time        `json:"created_date"`
}

type (
	PricePerUnit     domain.PricePerUnit
	MaterialType     domain.MaterialType
	MaterialQuantity domain.MaterialQuantity
)

type CropRead struct {
	UID        uuid.UUID  `json:"uid"`
	BatchID    string     `json:"batch_id"`
	Status     string     `json:"status"`
	Type       string     `json:"type"`
	Container  Container  `json:"container"`
	Inventory  Inventory  `json:"inventory"`
	AreaStatus AreaStatus `json:"area_status"`
	FarmUID    uuid.UUID  `json:"farm_id"`

	// Fields to track crop's movement
	InitialArea InitialArea `json:"initial_area"`
	MovedArea   []MovedArea `json:"moved_area"`
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

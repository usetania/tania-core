package storage

import (
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
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
	CountryCode string    `json:"country_code"`
	CityCode    string    `json:"city_code"`
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
	UID  uuid.UUID
	Name string
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

type AreaSize domain.AreaSize
type AreaLocation domain.AreaLocation
type AreaType domain.AreaType
type AreaPhoto domain.AreaPhoto
type AreaNote domain.AreaNote

type MaterialEvent struct {
	MaterialUID uuid.UUID
	Version     int
	Event       interface{}
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

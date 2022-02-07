package query

import (
	"time"

	"github.com/gofrs/uuid"
)

type FarmEvent interface {
	FindAllByID(farmUID uuid.UUID) <-chan Result
}

type FarmRead interface {
	FindByID(farmUID uuid.UUID) <-chan Result
	FindAll() <-chan Result
}

type ReservoirEvent interface {
	FindAllByID(reservoirUID uuid.UUID) <-chan Result
}

type ReservoirRead interface {
	FindByID(reservoirUID uuid.UUID) <-chan Result
	FindAllByFarm(farmUID uuid.UUID) <-chan Result
}

type AreaEvent interface {
	FindAllByID(areaUID uuid.UUID) <-chan Result
}

type AreaRead interface {
	FindByID(reservoirUID uuid.UUID) <-chan Result
	FindAllByFarm(farmUID uuid.UUID) <-chan Result
	FindByIDAndFarm(areaUID, farmUID uuid.UUID) <-chan Result
	FindAreasByReservoirID(reservoirUID uuid.UUID) <-chan Result
	CountAreas(farmUID uuid.UUID) <-chan Result
}

type CropRead interface {
	FindAllCropByArea(areaUID uuid.UUID) <-chan Result
	CountCropsByArea(areaUID uuid.UUID) <-chan Result
}

type MaterialEvent interface {
	FindAllByID(materialUID uuid.UUID) <-chan Result
}

type MaterialRead interface {
	FindAll(materialType, materialTypeDetail string, page, limit int) <-chan Result
	CountAll(materialType, materialTypeDetail string) <-chan Result
	FindByID(materialUID uuid.UUID) <-chan Result
}

type Result struct {
	Result interface{}
	Error  error
}

type FarmResult struct {
	UID         uuid.UUID
	Name        string
	Type        string
	Latitude    string
	Longitude   string
	CountryCode string
	CityCode    string
	CreatedDate time.Time
}

type ReservoirResult struct {
	UID         uuid.UUID
	Name        string
	WaterSource WaterSource
	FarmUID     uuid.UUID
	Notes       []ReservoirNote
	CreatedDate time.Time
}

type WaterSource struct {
	Type     string
	Capacity float32
}

type ReservoirNote struct {
	UID         uuid.UUID
	Content     string
	CreatedDate time.Time
}

type CountAreaCropResult struct {
	PlantQuantity  int
	TotalCropBatch int
}

type AreaCropResult struct {
	CropUID          uuid.UUID   `json:"uid"`
	BatchID          string      `json:"batch_id"`
	InitialArea      InitialArea `json:"initial_area"`
	MovingDate       time.Time   `json:"moving_date"`
	CreatedDate      time.Time   `json:"created_date"`
	DaysSinceSeeding int         `json:"days_since_seeding"`
	Inventory        Inventory   `json:"inventory"`
	Container        Container   `json:"container"`
}

type InitialArea struct {
	AreaUID uuid.UUID `json:"uid"`
	Name    string    `json:"name"`
}

type Inventory struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type Container struct {
	Type     ContainerType `json:"type"`
	Quantity int           `json:"quantity"`
}

type ContainerType struct {
	Code string `json:"code"`
	Cell int    `json:"cell"`
}

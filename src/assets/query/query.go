package query

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type FarmReadQuery interface {
	FindByID(farmUID uuid.UUID) <-chan QueryResult
	FindAll() <-chan QueryResult
}

type ReservoirEventQuery interface {
	FindAllByID(reservoirUID uuid.UUID) <-chan QueryResult
}

type ReservoirReadQuery interface {
	FindByID(reservoirUID uuid.UUID) <-chan QueryResult
	FindAll() <-chan QueryResult
}

type AreaQuery interface {
	FindAreasByReservoirID(reservoirID string) <-chan QueryResult
}

type CropQuery interface {
	FindAllCropByArea(areaUID uuid.UUID) <-chan QueryResult
	CountCropsByArea(areaUID uuid.UUID) <-chan QueryResult
}

type MaterialQuery interface {
	FindAll() <-chan QueryResult
	FindAllMaterialByPlantType(plantTypeCode string) <-chan QueryResult
	FindMaterialByPlantTypeAndName(plantTypeCode string, name string) <-chan QueryResult
}

type QueryResult struct {
	Result interface{}
	Error  error
}

type FarmReadQueryResult struct {
	UID         uuid.UUID
	Name        string
	Type        string
	Latitude    string
	Longitude   string
	CountryCode string
	CityCode    string
	CreatedDate time.Time
}

type ReservoirReadQueryResult struct {
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

type CountAreaCropQueryResult struct {
	PlantQuantity  int
	TotalCropBatch int
}

type AreaCropQueryResult struct {
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

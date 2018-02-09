package query

import (
	"time"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
)

type AreaQuery interface {
	FindAreasByReservoirID(reservoirID string) <-chan QueryResult
}

type CropQuery interface {
	FindAllCropByArea(areaUID uuid.UUID) <-chan QueryResult
	CountCropsByArea(areaUID uuid.UUID) <-chan QueryResult
}

type MaterialQuery interface {
	FindAll() <-chan QueryResult
	FindAllMaterialByPlantType(plantType domain.PlantType) <-chan QueryResult
	FindMaterialByPlantTypeAndName(plantType domain.PlantType, name string) <-chan QueryResult
}

type QueryResult struct {
	Result interface{}
	Error  error
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

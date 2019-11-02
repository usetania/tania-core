package query

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AreaReadQuery interface {
	FindByID(areaUID uuid.UUID) <-chan QueryResult
}

type CropQuery interface {
	FindByBatchID(batchID string) <-chan QueryResult
	FindAllCropsByFarm(farmUID uuid.UUID) <-chan QueryResult
	FindAllCropsByArea(areaUID uuid.UUID) <-chan QueryResult
}

type CropEventQuery interface {
	FindAllByCropID(uid uuid.UUID) <-chan QueryResult
}

type CropReadQuery interface {
	FindByID(uid uuid.UUID) <-chan QueryResult
	FindByBatchID(batchID string) <-chan QueryResult
	FindAllCropsByFarm(farmUID uuid.UUID, status string, page, limit int) <-chan QueryResult
	CountAllCropsByFarm(farmUID uuid.UUID, status string) <-chan QueryResult
	FindAllCropsByArea(areaUID uuid.UUID) <-chan QueryResult
	FindAllCropsArchives(farmUID uuid.UUID, page, limit int) <-chan QueryResult
	CountAllArchivedCropsByFarm(farmUID uuid.UUID) <-chan QueryResult
	FindCropsInformation(farmUID uuid.UUID) <-chan QueryResult
	CountTotalBatch(farmUID uuid.UUID) <-chan QueryResult
}

type CropActivityQuery interface {
	FindAllByCropID(uid uuid.UUID) <-chan QueryResult
	FindByCropIDAndActivityType(uid uuid.UUID, activityType interface{}) <-chan QueryResult
}

type MaterialReadQuery interface {
	FindByID(inventoryUID uuid.UUID) <-chan QueryResult
	FindMaterialByPlantTypeCodeAndName(plantType string, name string) <-chan QueryResult
}

type FarmReadQuery interface {
	FindByID(farmUID uuid.UUID) <-chan QueryResult
}

type TaskReadQuery interface {
	FindByID(taskUID uuid.UUID) <-chan QueryResult
}

type QueryResult struct {
	Result interface{}
	Error  error
}

type CropMaterialQueryResult struct {
	UID           uuid.UUID `json:"uid"`
	TypeCode      string    `json:"type"`
	PlantTypeCode string    `json:"plant_type"`
	Name          string    `json:"name"`
}

type CropAreaQueryResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
	Size struct {
		Value  float32 `json:"value"`
		Symbol string  `json:"symbol"`
	} `json:"size"`
	Type     string    `json:"type"`
	Location string    `json:"location"`
	FarmUID  uuid.UUID `json:"farm_uid"`
}

type CropAreaByAreaQueryResult struct {
	UID         uuid.UUID `json:"uid"`
	BatchID     string    `json:"batch_id"`
	Inventory   Inventory `json:"inventory"`
	CreatedDate time.Time `json:"seeding_date"`
	Area        Area      `json:"area"`
	Container   Container `json:"container"`
}

type Area struct {
	UID             uuid.UUID   `json:"uid"`
	Name            string      `json:"name"`
	InitialQuantity int         `json:"initial_quantity"`
	CurrentQuantity int         `json:"current_quantity"`
	InitialArea     InitialArea `json:"initial_area"`
	LastWatered     *time.Time  `json:"last_watered"`
	MovingDate      time.Time   `json:"moving_date"`
}

type InitialArea struct {
	UID         uuid.UUID `json:"uid"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

type Container struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Cell     int    `json:"cell"`
}

type Inventory struct {
	UID       uuid.UUID `json:"uid"`
	PlantType string    `json:"plant_type"`
	Name      string    `json:"name"`
}

type CropInformationQueryResult struct {
	TotalHarvestProduced float32 `json:"total_harvest_produced"`
	TotalPlantVariety    int     `json:"total_plant_variety"`
}

type CropFarmQueryResult struct {
	UID  uuid.UUID
	Name string
}

type CountTotalBatchQueryResult struct {
	VarietyName string `json:"variety_name"`
	TotalBatch  int    `json:"total_batch"`
}

type CropTaskQueryResult struct {
	UID         uuid.UUID
	Title       string
	Description string
	Category    string
	Status      string
	Domain      string
	AssetUID    uuid.UUID
	MaterialUID uuid.UUID
	AreaUID     uuid.UUID
}

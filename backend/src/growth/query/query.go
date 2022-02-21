package query

import (
	"time"

	"github.com/gofrs/uuid"
)

type AreaReadQuery interface {
	FindByID(areaUID uuid.UUID) <-chan Result
}

type CropQuery interface {
	FindByBatchID(batchID string) <-chan Result
	FindAllCropsByFarm(farmUID uuid.UUID) <-chan Result
	FindAllCropsByArea(areaUID uuid.UUID) <-chan Result
}

type CropEventQuery interface {
	FindAllByCropID(uid uuid.UUID) <-chan Result
}

type CropReadQuery interface {
	FindByID(uid uuid.UUID) <-chan Result
	FindByBatchID(batchID string) <-chan Result
	FindAllCropsByFarm(farmUID uuid.UUID, status string, page, limit int) <-chan Result
	CountAllCropsByFarm(farmUID uuid.UUID, status string) <-chan Result
	FindAllCropsByArea(areaUID uuid.UUID) <-chan Result
	FindAllCropsArchives(farmUID uuid.UUID, page, limit int) <-chan Result
	CountAllArchivedCropsByFarm(farmUID uuid.UUID) <-chan Result
	FindCropsInformation(farmUID uuid.UUID) <-chan Result
	CountTotalBatch(farmUID uuid.UUID) <-chan Result
}

type CropActivityQuery interface {
	FindAllByCropID(uid uuid.UUID) <-chan Result
	FindByCropIDAndActivityType(uid uuid.UUID, activityType interface{}) <-chan Result
}

type MaterialReadQuery interface {
	FindByID(inventoryUID uuid.UUID) <-chan Result
	FindMaterialByPlantTypeCodeAndName(plantType string, name string) <-chan Result
}

type FarmReadQuery interface {
	FindByID(farmUID uuid.UUID) <-chan Result
}

type TaskReadQuery interface {
	FindByID(taskUID uuid.UUID) <-chan Result
}

type Result struct {
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

package query

import uuid "github.com/satori/go.uuid"

type AreaQuery interface {
	FindByID(areaUID uuid.UUID) <-chan QueryResult
}

type CropQuery interface {
	FindByBatchID(batchID string) <-chan QueryResult
	FindAllCropsByFarm(farmUID uuid.UUID) <-chan QueryResult
	FindAllCropsByArea(areaUID uuid.UUID) <-chan QueryResult
}

type InventoryMaterialQuery interface {
	FindByID(inventoryUID uuid.UUID) <-chan QueryResult
	FindInventoryByPlantTypeCodeAndVariety(plantType string, variety string) <-chan QueryResult
}

type FarmQuery interface {
	FindByID(farmUID uuid.UUID) <-chan QueryResult
}

type QueryResult struct {
	Result interface{}
	Error  error
}

type CropInventoryQueryResult struct {
	UID           uuid.UUID `json:"uid"`
	PlantTypeCode string    `json:"plant_type"`
	Variety       string    `json:"variety"`
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

type CropFarmQueryResult struct {
	UID  uuid.UUID
	Name string
}

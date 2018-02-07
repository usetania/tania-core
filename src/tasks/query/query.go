package query

import uuid "github.com/satori/go.uuid"

type QueryResult struct {
	Result interface{}
	Error  error
}

type AreaQuery interface {
	FindByID(areaUID uuid.UUID) <-chan QueryResult
}

type CropQuery interface {
	FindCropByID(cropUID uuid.UUID) <-chan QueryResult
}

type MaterialQuery interface {
	FindMaterialByID(materialID uuid.UUID) <-chan QueryResult
}

type TaskQuery interface {
	QueryTasksWithFilter(params map[string]string) <-chan QueryResult
}

/*
TODO

type ReservoirQuery interface {
	FindReservoiryByID(reservoirUID uuid.UUID) <-chan QueryResult
}

type DeviceQuery interface {
	FindDeviceByID(deviceUID uuid.UUID) <-chan QueryResult
}

type FinanceQuery interface {
	FindFinanceByID(financeUID uuid.UUID) <-chan QueryResult
}

*/
type TaskAreaQueryResult struct {
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

type TaskCropQueryResult struct {
	UID       uuid.UUID `json:"uid"`
	BatchID   string    `json:"batch_id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Container struct {
		Quantity int    `json:"quantity"`
		Type     string `json:"type"`
	} `json:"container"`
	InventoryUID uuid.UUID `json:"inventory_uid"`
	FarmUID      uuid.UUID `json:"farm_uid"`
}

type TaskMaterialQueryResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

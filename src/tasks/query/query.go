package query

import (
	uuid "github.com/satori/go.uuid"
)

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

type TaskEventQuery interface {
	FindAllByTaskID(uid uuid.UUID) <-chan QueryResult
}

type TaskReadQuery interface {
	FindAll() <-chan QueryResult
	FindByID(string) <-chan QueryResult
	QueryTasksWithFilter(params map[string]string) <-chan QueryResult
}

type ReservoirQuery interface {
	FindReservoirByID(reservoirUID uuid.UUID) <-chan QueryResult
}

/*
TODO

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
}

type TaskCropQueryResult struct {
	UID     uuid.UUID `json:"uid"`
	BatchID string    `json:"batch_id"`
}

type TaskMaterialQueryResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type TaskReservoirQueryResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

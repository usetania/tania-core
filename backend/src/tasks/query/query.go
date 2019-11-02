package query

import (
	//assetsdomain "github.com/Tanibox/tania-core/src/assets/domain"
	uuid "github.com/satori/go.uuid"
)

type QueryResult struct {
	Result interface{}
	Error  error
}

// EventWrapper is used to wrap the event interface with its struct name,
// so it will be easier to unmarshal later
type EventWrapper struct {
	EventName string
	EventData interface{}
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

type TaskEventQuery interface {
	FindAllByTaskID(uid uuid.UUID) <-chan QueryResult
}

type TaskReadQuery interface {
	FindAll(page, limit int) <-chan QueryResult
	FindByID(taskUID uuid.UUID) <-chan QueryResult
	FindTasksWithFilter(params map[string]string, page, limit int) <-chan QueryResult
  CountAll() <-chan QueryResult
  CountTasksWithFilter(params map[string]string) <-chan QueryResult
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

// QUERY RESULTS

type TaskAreaQueryResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type TaskCropQueryResult struct {
	UID     uuid.UUID `json:"uid"`
	BatchID string    `json:"batch_id"`
}

type TaskMaterialQueryResult struct {
	UID              uuid.UUID `json:"uid"`
	TypeCode         string    `json:"type"`
	DetailedTypeCode string    `json:"detailed_type"`
	Name             string    `json:"name"`
}

type TaskReservoirQueryResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

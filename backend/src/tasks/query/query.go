package query

import (
	"github.com/gofrs/uuid"
)

type Result struct {
	Result interface{}
	Error  error
}

// EventWrapper is used to wrap the event interface with its struct name,
// so it will be easier to unmarshal later.
type EventWrapper struct {
	EventName string
	EventData interface{}
}

type Area interface {
	FindByID(areaUID uuid.UUID) <-chan Result
}

type Crop interface {
	FindCropByID(cropUID uuid.UUID) <-chan Result
}

type Material interface {
	FindMaterialByID(materialID uuid.UUID) <-chan Result
}

type TaskEvent interface {
	FindAllByTaskID(uid uuid.UUID) <-chan Result
}

type TaskRead interface {
	FindAll(page, limit int) <-chan Result
	FindByID(taskUID uuid.UUID) <-chan Result
	FindTasksWithFilter(params map[string]string, page, limit int) <-chan Result
	CountAll() <-chan Result
	CountTasksWithFilter(params map[string]string) <-chan Result
}

type Reservoir interface {
	FindReservoirByID(reservoirUID uuid.UUID) <-chan Result
}

// QUERY RESULTS

type TaskAreaResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

type TaskCropResult struct {
	UID     uuid.UUID `json:"uid"`
	BatchID string    `json:"batch_id"`
}

type TaskMaterialResult struct {
	UID              uuid.UUID `json:"uid"`
	TypeCode         string    `json:"type"`
	DetailedTypeCode string    `json:"detailed_type"`
	Name             string    `json:"name"`
}

type TaskReservoirResult struct {
	UID  uuid.UUID `json:"uid"`
	Name string    `json:"name"`
}

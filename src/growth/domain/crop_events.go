package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CropBatchCreated struct {
	UID            uuid.UUID
	BatchID        string
	Status         CropStatus
	Type           CropType
	Container      CropContainer
	InventoryUID   uuid.UUID
	FarmUID        uuid.UUID
	CreatedDate    time.Time
	InitialAreaUID uuid.UUID
	Quantity       int
}

type CropBatchMoved struct {
	UID            uuid.UUID
	Quantity       int
	SrcAreaUID     uuid.UUID
	DstAreaUID     uuid.UUID
	MovedDate      time.Time
	UpdatedSrcArea interface{}
	UpdatedDstArea MovedArea
}

type CropBatchHarvested struct {
	UID                     uuid.UUID
	HarvestedQuantity       int
	ProducedGramQuantity    float32
	UpdatedHarvestedStorage HarvestedStorage
	HarvestedArea           interface{}
	HarvestedAreaType       string
	HarvestDate             time.Time
}

type CropBatchDumped struct {
	UID            uuid.UUID
	Quantity       int
	UpdatedTrash   Trash
	DumpedArea     interface{}
	DumpedAreaType string
	DumpDate       time.Time
}

type CropBatchWatered struct {
	UID           uuid.UUID
	BatchID       string
	ContainerType string
	AreaUID       uuid.UUID
	AreaName      string
	WateringDate  time.Time
}

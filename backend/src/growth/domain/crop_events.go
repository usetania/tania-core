package domain

import (
	"time"

	"github.com/gofrs/uuid"
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

type CropBatchTypeChanged struct {
	UID  uuid.UUID
	Type CropType
}

type CropBatchInventoryChanged struct {
	UID          uuid.UUID
	InventoryUID uuid.UUID
	BatchID      string
}

type CropBatchContainerChanged struct {
	UID       uuid.UUID
	Container CropContainer
}

type CropBatchMoved struct {
	UID                uuid.UUID
	Quantity           int
	SrcAreaUID         uuid.UUID
	DstAreaUID         uuid.UUID
	MovedDate          time.Time
	UpdatedSrcAreaCode string // Values: INITIAL_AREA / MOVED_AREA
	UpdatedSrcArea     interface{}
	UpdatedDstAreaCode string // Values: INITIAL_AREA / MOVED_AREA
	UpdatedDstArea     interface{}
}

type CropBatchHarvested struct {
	UID                     uuid.UUID
	CropStatus              string // Values: ACTIVE / ARCHIVED
	HarvestType             string
	HarvestedQuantity       int
	ProducedGramQuantity    float32
	UpdatedHarvestedStorage HarvestedStorage
	HarvestedArea           interface{}
	HarvestedAreaCode       string // Values: INITIAL_AREA / MOVED_AREA
	HarvestDate             time.Time
	Notes                   string
}

type CropBatchDumped struct {
	UID            uuid.UUID
	CropStatus     string // Values: ACTIVE / ARCHIVED
	Quantity       int
	UpdatedTrash   Trash
	DumpedArea     interface{}
	DumpedAreaCode string // Values: INITIAL_AREA / MOVED_AREA
	DumpDate       time.Time
	Notes          string
}

type CropBatchWatered struct {
	UID           uuid.UUID
	BatchID       string
	ContainerType string
	AreaUID       uuid.UUID
	AreaName      string
	WateringDate  time.Time
}

type CropBatchNoteCreated struct {
	UID         uuid.UUID
	CropUID     uuid.UUID
	Content     string
	CreatedDate time.Time
}

type CropBatchNoteRemoved struct {
	UID         uuid.UUID
	CropUID     uuid.UUID
	Content     string
	CreatedDate time.Time
}

type CropBatchPhotoCreated struct {
	UID         uuid.UUID
	CropUID     uuid.UUID
	Filename    string
	MimeType    string
	Size        int
	Width       int
	Height      int
	Description string
}

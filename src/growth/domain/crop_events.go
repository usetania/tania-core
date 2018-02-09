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
	UpdatedDstArea interface{}
}

type CropBatchHarvested struct {
	UID                     uuid.UUID
	HarvestedQuantity       int
	ProducedGramQuantity    float32
	UpdatedHarvestedStorage HarvestedStorage
	HarvestedArea           interface{}
	HarvestedAreaType       string
	HarvestDate             time.Time
	Notes                   string
}

type CropBatchDumped struct {
	UID            uuid.UUID
	Quantity       int
	UpdatedTrash   Trash
	DumpedArea     interface{}
	DumpedAreaType string
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

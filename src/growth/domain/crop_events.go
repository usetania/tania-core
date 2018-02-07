package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CropBatchCreated struct {
	UID             uuid.UUID
	BatchID         string
	Status          CropStatus
	Type            CropType
	Container       CropContainer
	ContainerType   string
	ContainerCell   int
	InventoryUID    uuid.UUID
	VarietyName     string
	PlantType       string
	FarmUID         uuid.UUID
	CreatedDate     time.Time
	InitialAreaUID  uuid.UUID
	InitialAreaName string
	Quantity        int
}

type CropBatchMoved struct {
	UID                    uuid.UUID
	BatchID                string
	ContainerType          string
	Quantity               int
	SrcAreaUID             uuid.UUID
	SrcAreaName            string
	SrcAreaCurrentQuantity int
	SrcAreaType            string
	DstAreaUID             uuid.UUID
	DstAreaName            string
	DstAreaCurrentQuantity int
	DstAreaType            string
	MovedDate              time.Time
}

type CropBatchHarvested struct {
	UID                  uuid.UUID
	BatchID              string
	ContainerType        string
	HarvestType          string
	Quantity             int
	ProducedGramQuantity float32
	SrcAreaUID           uuid.UUID
	SrcAreaName          string
	HarvestDate          time.Time
}

type CropBatchDumped struct {
	UID           uuid.UUID
	BatchID       string
	ContainerType string
	Quantity      int
	SrcAreaUID    uuid.UUID
	SrcAreaName   string
	SrcAreaType   string
	DumpDate      time.Time
}

type CropBatchWatered struct {
	UID           uuid.UUID
	BatchID       string
	ContainerType string
	AreaUID       uuid.UUID
	AreaName      string
	WateringDate  time.Time
}

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
	FarmUID         uuid.UUID
	CreatedDate     time.Time
	InitialAreaUID  uuid.UUID
	InitialAreaName string
	Quantity        int
}

type CropBatchWatered struct {
	UID          uuid.UUID
	AreaUID      uuid.UUID
	AreaName     string
	WateringDate time.Time
}

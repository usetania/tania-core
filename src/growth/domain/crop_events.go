package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CropBatchCreated struct {
	UID          uuid.UUID
	BatchID      string
	Status       CropStatus
	Type         CropType
	Container    CropContainer
	InventoryUID uuid.UUID
	FarmUID      uuid.UUID
	CreatedDate  time.Time
	Photo        CropPhoto

	// Fields to track crop's movement
	InitialArea InitialArea
}

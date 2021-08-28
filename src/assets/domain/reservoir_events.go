package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type ReservoirCreated struct {
	UID         uuid.UUID
	Name        string
	WaterSource WaterSource
	FarmUID     uuid.UUID
	CreatedDate time.Time
}

type ReservoirWaterSourceChanged struct {
	ReservoirUID uuid.UUID
	WaterSource  WaterSource
}

type ReservoirNameChanged struct {
	ReservoirUID uuid.UUID
	Name         string
}

type ReservoirNoteAdded struct {
	ReservoirUID uuid.UUID
	UID          uuid.UUID
	Content      string
	CreatedDate  time.Time
}

type ReservoirNoteRemoved struct {
	ReservoirUID uuid.UUID
	UID          uuid.UUID
}

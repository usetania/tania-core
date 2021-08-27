package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type AreaCreated struct {
	UID          uuid.UUID
	Name         string
	Type         AreaType
	Location     AreaLocation
	Size         AreaSize
	FarmUID      uuid.UUID
	ReservoirUID uuid.UUID
	CreatedDate  time.Time
}

type AreaNameChanged struct {
	AreaUID uuid.UUID
	Name    string
}

type AreaSizeChanged struct {
	AreaUID uuid.UUID
	Size    AreaSize
}

type AreaTypeChanged struct {
	AreaUID uuid.UUID
	Type    AreaType
}

type AreaLocationChanged struct {
	AreaUID  uuid.UUID
	Location AreaLocation
}

type AreaReservoirChanged struct {
	AreaUID      uuid.UUID
	ReservoirUID uuid.UUID
}

type AreaPhotoAdded struct {
	AreaUID  uuid.UUID
	Filename string
	MimeType string
	Size     int
	Width    int
	Height   int
}

type AreaNoteAdded struct {
	AreaUID     uuid.UUID
	UID         uuid.UUID
	Content     string
	CreatedDate time.Time
}

type AreaNoteRemoved struct {
	AreaUID uuid.UUID
	UID     uuid.UUID
}

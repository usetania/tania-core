package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type FarmCreated struct {
	UID         uuid.UUID
	UserID      uuid.UUID
	Name        string
	Type        string
	Latitude    string
	Longitude   string
	Country     string
	City        string
	IsActive    bool
	CreatedDate time.Time
}

type FarmNameChanged struct {
	FarmUID uuid.UUID
	Name    string
}

type FarmTypeChanged struct {
	FarmUID uuid.UUID
	Type    string
}

type FarmGeolocationChanged struct {
	FarmUID   uuid.UUID
	Latitude  string
	Longitude string
}

type FarmRegionChanged struct {
	FarmUID uuid.UUID
	Country string
	City    string
}

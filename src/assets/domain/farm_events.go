package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type FarmCreated struct {
	UID         uuid.UUID
	Name        string
	Type        string
	Latitude    string
	Longitude   string
	CountryCode string
	CityCode    string
	IsActive    bool
	CreatedDate time.Time
}

type FarmGeolocationChanged struct {
	Latitude  string
	Longitude string
}

type FarmRegionChanged struct {
	CountryCode string
	CityCode    string
}

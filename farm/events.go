// Package farm provides the operation that farm holder can do
// to their farm
package farm

type FarmCreated struct {
	Uid         string
	Name        string
	Description string
	Latitude    float64
	Longitude   float64
}

type FarmInformationUpdated struct {
	NameUpdated        string
	DescriptionUpdated string
	LatitudeUpdated    float64
	LongitudeUpdated   float64
}

type FarmDestroyed struct{}

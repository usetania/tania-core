// Package farm provides the operation that farm holder can do
// to their farm
package farm

type Farm struct {
	uid         string
	name        string
	description string
	latitude    float64
	longitude   float64
}

type Operation interface {
	// CreateFarm registers a new farm to Tania
	CreateFarm()

	// UpdateFarmInformation updates the existing farm information in Tania
	UpdateFarmInformation()

	// DestroyFarm destroys the farm and its properties. This is dangerous.
	DestroyFarm()
}

type EventHandle interface {
	FarmCreated()
	FarmInformationUpdated()
	FarmDestroyed()
}

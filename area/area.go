// Package area provides the operation that farm owner and his/her workers staff
// can do on the area of the farm
package area

type Operation interface {
	// CreateArea registers a new area inside the farm
	CreateArea()

	// HarvestCrops harvests all crops in an area
	HarvestCrops()

	// DisposeCrops disposes all crops in an area
	DisposeCrops()

	// IrrigateArea irrigates all crops in an area
	IrrigateArea()

	// FertilizeArea fertilizes all crops in an area
	FertilizeArea()

	// DestroyArea destroys an area and its plant on it.
	DestroyArea()
}

type EventHandle interface {
	AreaCreated()
	CropsHarvested()
	CropsDisposed()
	AreaIrrigated()
	AreaFertilized()
	AreaDestroyed()
}

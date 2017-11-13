// Package area provides the operation that farm owner and his/her workers staff
// can do on the area of the farm
package area

type EventHandler interface {
	areaCreated()
	cropsHarvested()
	cropsDisposed()
	areaIrrigated()
	areaFertilized()
	areaDestroyed()
}

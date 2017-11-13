// Package area provides the operation that farm owner and his/her workers staff
// can do on the area of the farm
package area

type AreaCreated struct {
	Uid           string
	GrowingMethod string
	Size          float64
}

type CropsHarvested struct{}

type CropsDisposed struct {
	ReasonAdded string
}

type AreaDestroyed struct{}

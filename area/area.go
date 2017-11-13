// Package area provides the operation that farm owner and his/her workers staff
// can do on the area of the farm
package area

type Area struct {
	Uid           string
	GrowingMethod string
	Size          float64
}

// DisplayAll shows all existing areas inside a farm
func DisplayAll() []Area {
	var areas []Area

	for i := 0; i < 5; i++ {
		areas = append(areas, Area{
			Uid:           string(i),
			GrowingMethod: "Dummy growing method",
			Size:          50.5,
		})
	}

	return areas
}

// CreateNew registers a new area inside the farm
func CreateNew(uid int) {
}

// HarvestCrops harvests all crops in an area
func HarvestCrops(uid int) {
}

// DisposeCrops disposes all crops in an area
func DisposeCrops(uid int) {
}

// Destroy destroys an area and its plant on it.
func Destroy(uid int) {
}

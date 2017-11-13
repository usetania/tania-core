// Package farm provides the operation that farm holder can do
// to their farm
package farm

type Farm struct {
	Uid         string
	Name        string
	Description string
	Latitude    float64
	Longitude   float64
}

// DisplayAll dispalys all existing farms
func DisplayAll() []Farm {
	var information []Farm

	for i := 0; i < 5; i++ {
		information = append(information, Farm{
			Uid:         string(i),
			Name:        "Dummy farm",
			Description: "",
			Latitude:    1.335678,
			Longitude:   5.34352,
		})
	}

	return information
}

// CreateNew registers a new farm to Tania
func CreateNew() {
}

// ShowInformation shows information of a farm
func ShowInformation(uid string) *Farm {
	information := &Farm{
		Uid:         uid,
		Name:        "Dummy farm",
		Description: "",
		Latitude:    1.335678,
		Longitude:   5.34352,
	}

	return information
}

// UpdateInformation updates the existing farm information in Tania
func UpdateInformation() {
}

// DestroyFarm destroys the farm and its properties. This is dangerous.
func Destroy() {

}

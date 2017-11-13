// Package reservoir provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm
package reservoir

type Reservoir struct {
	Uid         string
	Name        string
	Description string
	Capacity    float64
}

// DisplayAll shows all reservoirs in a farm
func DisplayAll() []Reservoir {
	var reservoirs []Reservoir

	for i := 0; i < 5; i++ {
		reservoirs = append(reservoirs, Reservoir{
			Uid:         string(i),
			Name:        "Dummy reservoir",
			Description: "Just a description",
			Capacity:    160.58,
		})
	}

	return reservoirs
}

// CreateNew registers a new reservoir in a farm
func CreateNew() {
}

// UpdateInformation updates information of the particular reservoir
func UpdateInformation(uid string) {
}

// Destroy destroys the reservoir. This is dangerous.
func Destroy(uid string) {
}

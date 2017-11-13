// Package reservoir provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm
package reservoir

type ReservoirCreated struct {
	Uid         string
	Name        string
	Description string
	Capacity    float64
}

type ReservoirInformationUpdated struct {
	NameUpdated        string
	DescriptionUpdated string
	CapacityUpdated    float64
}

type ReservoirDestroyed struct{}

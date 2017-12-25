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

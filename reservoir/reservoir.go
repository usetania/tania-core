package reservoir

type Reservoir struct {
	uid         string
	name        string
	description string
	capacity    float64
	createdAt   string
	updatedAt   string
}

type Operation interface {
	// createNew registers a new reservoir in the farm
	createNew()

	// updateInformation updates information of the reservoir
	updateInformation()

	// destroy destroys the reservoir. This is dangerous.
	destroy()
}

type EventHandler interface {
	reservoirCreated()
	reservoirInformationUpdated()
	reservoirDestroyed()
}

// Package reservoir provides the operation that farm owner or his/her staff
// can do with the reservoir in a farm
package reservoir

type EventHandler interface {
	reservoirCreated()
	reservoirInformationUpdated()
	reservoirDestroyed()
}

// Package farm provides the operation that farm holder can do
// to their farm
package farm

type EventHandler interface {
	farmCreated()
	farmInformationUpdated()
	farmDestroyed()
}

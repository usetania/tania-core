package repository

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/storage"
)

// Result is a struct to wrap repository result
// so its easy to use it in channel.
type Result struct {
	Result interface{}
	Error  error
}

type FarmEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type FarmRead interface {
	Save(farmRead *storage.FarmRead) <-chan error
}

func NewFarmFromHistory(events []storage.FarmEvent) *domain.Farm {
	state := &domain.Farm{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

type AreaEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type AreaRead interface {
	Save(areaRead *storage.AreaRead) <-chan error
}

func NewAreaFromHistory(events []storage.AreaEvent) *domain.Area {
	state := &domain.Area{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

type ReservoirEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type ReservoirRead interface {
	Save(reservoirRead *storage.ReservoirRead) <-chan error
}

func NewReservoirFromHistory(events []storage.ReservoirEvent) *domain.Reservoir {
	state := &domain.Reservoir{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

type MaterialEvent interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

func NewMaterialFromHistory(events []storage.MaterialEvent) *domain.Material {
	state := &domain.Material{}
	for _, v := range events {
		state.Transition(v.Event)
		state.Version++
	}

	return state
}

type MaterialEventTypeWrapper struct {
	Type string
	Data interface{}
}

func (w MaterialEventTypeWrapper) Code() string {
	return w.Type
}

type MaterialRead interface {
	Save(materialRead *storage.MaterialRead) <-chan error
}

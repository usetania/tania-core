package repository

import (
	"github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Result interface{}
	Error  error
}

type FarmEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type FarmReadRepository interface {
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

type AreaEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type AreaReadRepository interface {
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

type ReservoirEventRepository interface {
	Save(uid uuid.UUID, latestVersion int, events []interface{}) <-chan error
}

type ReservoirReadRepository interface {
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

type MaterialEventRepository interface {
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

type MaterialReadRepository interface {
	Save(materialRead *storage.MaterialRead) <-chan error
}

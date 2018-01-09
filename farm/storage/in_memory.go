package storage

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
	uuid "github.com/satori/go.uuid"
)

type FarmStorage struct {
	Lock    sync.RWMutex
	FarmMap map[uuid.UUID]entity.Farm
}

type AreaStorage struct {
	Lock    sync.RWMutex
	AreaMap map[uuid.UUID]entity.Area
}

type ReservoirStorage struct {
	Lock         sync.RWMutex
	ReservoirMap map[uuid.UUID]entity.Reservoir
}

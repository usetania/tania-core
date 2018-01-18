package storage

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type FarmStorage struct {
	Lock    *deadlock.RWMutex
	FarmMap map[uuid.UUID]domain.Farm
}

type AreaStorage struct {
	Lock    *deadlock.RWMutex
	AreaMap map[uuid.UUID]domain.Area
}

type ReservoirStorage struct {
	Lock         *deadlock.RWMutex
	ReservoirMap map[uuid.UUID]domain.Reservoir
}

type InventoryMaterialStorage struct {
	Lock                 *deadlock.RWMutex
	InventoryMaterialMap map[uuid.UUID]domain.InventoryMaterial
}

type CropStorage struct {
	Lock    *deadlock.RWMutex
	CropMap map[uuid.UUID]domain.Crop
}

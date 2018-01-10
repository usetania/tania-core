package storage

import (
	"sync"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
)

type FarmStorage struct {
	Lock    sync.RWMutex
	FarmMap map[uuid.UUID]domain.Farm
}

type AreaStorage struct {
	Lock    sync.RWMutex
	AreaMap map[uuid.UUID]domain.Area
}

type ReservoirStorage struct {
	Lock         sync.RWMutex
	ReservoirMap map[uuid.UUID]domain.Reservoir
}

type InventoryMaterialStorage struct {
	Lock                 sync.RWMutex
	InventoryMaterialMap map[uuid.UUID]domain.InventoryMaterial
}

type CropStorage struct {
	Lock    sync.RWMutex
	CropMap map[uuid.UUID]domain.Crop
}

package storage

import (
	"sync"

	"github.com/Tanibox/tania-server/farm/entity"
	uuid "github.com/satori/go.uuid"
)

type AreaStorage struct {
	Lock    sync.RWMutex
	AreaMap map[uuid.UUID]entity.Area
}

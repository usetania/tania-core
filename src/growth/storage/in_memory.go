package storage

import (
	"github.com/Tanibox/tania-server/src/growth/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type CropStorage struct {
	Lock    *deadlock.RWMutex
	CropMap map[uuid.UUID]domain.Crop
}

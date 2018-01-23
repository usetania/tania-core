package storage

import (
	"fmt"
	"time"

	"github.com/Tanibox/tania-server/src/growth/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type CropStorage struct {
	Lock    *deadlock.RWMutex
	CropMap map[uuid.UUID]domain.Crop
}

func CreateCropStorage() *CropStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("DEADLOCK!")
	}

	return &CropStorage{CropMap: make(map[uuid.UUID]domain.Crop), Lock: &rwMutex}
}

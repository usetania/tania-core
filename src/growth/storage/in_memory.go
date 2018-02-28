package storage

import (
	"fmt"
	"time"

	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
)

type CropEventStorage struct {
	Lock       *deadlock.RWMutex
	CropEvents []CropEvent
}

type CropReadStorage struct {
	Lock        *deadlock.RWMutex
	CropReadMap map[uuid.UUID]CropRead
}

func CreateCropReadStorage() *CropReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP READ STORAGE DEADLOCK!")
	}

	return &CropReadStorage{CropReadMap: make(map[uuid.UUID]CropRead), Lock: &rwMutex}
}

type CropActivityStorage struct {
	Lock            *deadlock.RWMutex
	CropActivityMap []CropActivity
}

func CreateCropActivityStorage() *CropActivityStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("CROP LIST STORAGE DEADLOCK!")
	}

	return &CropActivityStorage{CropActivityMap: []CropActivity{}, Lock: &rwMutex}
}

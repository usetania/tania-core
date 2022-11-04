package storage

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sasha-s/go-deadlock"
)

type TaskEventStorage struct {
	Lock       *deadlock.RWMutex
	TaskEvents []TaskEvent
}

func CreateTaskEventStorage() *TaskEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		log.Println("TASK EVENT STORAGE DEADLOCK!")
	}

	return &TaskEventStorage{Lock: &rwMutex}
}

type TaskReadStorage struct {
	Lock        *deadlock.RWMutex
	TaskReadMap map[uuid.UUID]TaskRead
}

func CreateTaskReadStorage() *TaskReadStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		log.Println("TASK READ STORAGE DEADLOCK!")
	}

	return &TaskReadStorage{TaskReadMap: make(map[uuid.UUID]TaskRead), Lock: &rwMutex}
}

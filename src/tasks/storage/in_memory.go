package storage

import (
	"fmt"
	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	deadlock "github.com/sasha-s/go-deadlock"
	uuid "github.com/satori/go.uuid"
	"time"
)

type TaskStorage struct {
	Lock    *deadlock.RWMutex
	TaskMap map[uuid.UUID]domain.Task
}

type TaskEventStorage struct {
	Lock       *deadlock.RWMutex
	TaskEvents []TaskEvent
}

type TaskEvent struct {
	TaskUID uuid.UUID
	Version int
	Event   interface{}
}

func CreateTaskEventStorage() *TaskEventStorage {
	rwMutex := deadlock.RWMutex{}
	deadlock.Opts.DeadlockTimeout = time.Second * 10
	deadlock.Opts.OnPotentialDeadlock = func() {
		fmt.Println("TASK EVENT STORAGE DEADLOCK!")
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
		fmt.Println("TASK READ STORAGE DEADLOCK!")
	}

	return &TaskReadStorage{TaskReadMap: make(map[uuid.UUID]TaskRead), Lock: &rwMutex}
}

type TaskRead struct {
	Title         string            `json:"title"`
	UID           uuid.UUID         `json:"uid"`
	Description   string            `json:"description"`
	CreatedDate   time.Time         `json:"created_date"`
	DueDate       *time.Time        `json:"due_date, omitempty"`
	CompletedDate *time.Time        `json:"completed_date"`
	Priority      string            `json:"priority"`
	Status        string            `json:"status"`
	Domain        string            `json:"domain"`
	DomainDetails domain.TaskDomain `json:"domain_details"`
	Category      string            `json:"category"`
	IsDue         bool              `json:"is_due"`
	AssetID       *uuid.UUID        `json:"asset_id"`
}

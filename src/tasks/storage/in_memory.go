package storage

import (
	"sync"
	domain "github.com/Tanibox/tania-server/src/tasks/domain"

	uuid "github.com/satori/go.uuid"
)

type TaskStorage struct {
	Lock    sync.RWMutex
	TaskMap map[uuid.UUID]domain.Task
}

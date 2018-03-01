package storage

import (
	"time"

	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	uuid "github.com/satori/go.uuid"
)

type TaskEvent struct {
	TaskUID     uuid.UUID
	Version     int
	CreatedDate time.Time
	Event       interface{}
}

type TaskRead struct {
	Title         string            `json:"title"`
	UID           uuid.UUID         `json:"uid"`
	Description   string            `json:"description"`
	CreatedDate   time.Time         `json:"created_date"`
	DueDate       *time.Time        `json:"due_date, omitempty"`
	CompletedDate *time.Time        `json:"completed_date"`
	CancelledDate *time.Time        `json:"cancelled_date"`
	Priority      string            `json:"priority"`
	Status        string            `json:"status"`
	Domain        string            `json:"domain"`
	DomainDetails domain.TaskDomain `json:"domain_details"`
	Category      string            `json:"category"`
	IsDue         bool              `json:"is_due"`
	AssetID       *uuid.UUID        `json:"asset_id"`
}

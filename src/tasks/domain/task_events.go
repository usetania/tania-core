package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	TaskCreatedCode            = "TaskCreated"
	TaskTitleChangedCode       = "TaskTitleChanged"
	TaskDescriptionChangedCode = "TaskDescriptionChanged"
	TaskPriorityChangedCode    = "TaskPriorityChanged"
	TaskDueDateChangedCode     = "TaskDueDateChanged"
	TaskCategoryChangedCode    = "TaskCategoryChanged"
	TaskDetailsChangedCode     = "TaskDetailsChanged"
	TaskAssetIDChangedCode     = "TaskAssetIDChanged"
	TaskCompletedCode          = "TaskCompleted"
	TaskCancelledCode          = "TaskCancelled"
	TaskDueCode                = "TaskDue"
)

type TaskCreated struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CreatedDate   time.Time  `json:"created_date"`
	DueDate       *time.Time `json:"due_date"`
	Priority      string     `json:"priority"`
	Status        string     `json:"status"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	IsDue         bool       `json:"is_due"`
	AssetID       *uuid.UUID `json:"asset_id"`
}

type TaskTitleChanged struct {
	UID   uuid.UUID `json:"uid"`
	Title string    `json:"title"`
}

type TaskDescriptionChanged struct {
	UID         uuid.UUID `json:"uid"`
	Description string    `json:"description"`
}

type TaskPriorityChanged struct {
	UID      uuid.UUID `json:"uid"`
	Priority string    `json:"priority"`
}

type TaskDueDateChanged struct {
	UID     uuid.UUID  `json:"uid"`
	DueDate *time.Time `json:"due_date"`
}

type TaskCategoryChanged struct {
	UID      uuid.UUID `json:"uid"`
	Category string    `json:"category"`
}

type TaskDetailsChanged struct {
	UID           uuid.UUID  `json:"uid"`
	DomainDetails TaskDomain `json:"domain_details"`
}

type TaskAssetIDChanged struct {
	UID     uuid.UUID  `json:"uid"`
	AssetID *uuid.UUID `json:"asset_id"`
}

type TaskCompleted struct {
	UID           uuid.UUID  `json:"uid"`
	Status        string     `json:"status"`
	CompletedDate *time.Time `json:"completed_date"`
}

type TaskCancelled struct {
	UID           uuid.UUID  `json:"uid"`
	Status        string     `json:"status"`
	CancelledDate *time.Time `json:"cancelled_date"`
}

type TaskDue struct {
	UID uuid.UUID `json:"uid"`
}

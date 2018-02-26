package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	TaskCreatedCode   = "TaskCreated"
	TaskCompletedCode = "TaskCompleted"
	TaskCancelledCode = "TaskCancelled"
	TaskDueCode       = "TaskDue"
	TaskModifiedCode  = "TaskModified"
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

type TaskModified struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Priority      string     `json:"priority"`
	DueDate       *time.Time `json:"due_date"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	AssetID       *uuid.UUID `json:"asset_id"`
}

type TaskCompleted struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Priority      string     `json:"priority"`
	DueDate       *time.Time `json:"due_date"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	AssetID       *uuid.UUID `json:"asset_id"`
	CompletedDate *time.Time `json:"completed_date"`
}

type TaskCancelled struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Priority      string     `json:"priority"`
	DueDate       *time.Time `json:"due_date"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	AssetID       *uuid.UUID `json:"asset_id"`
	CancelledDate *time.Time `json:"cancelled_date"`
}

type TaskDue struct {
	UID           uuid.UUID  `json:"uid"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Priority      string     `json:"priority"`
	DueDate       *time.Time `json:"due_date"`
	Domain        string     `json:"domain"`
	DomainDetails TaskDomain `json:"domain_details"`
	Category      string     `json:"category"`
	AssetID       *uuid.UUID `json:"asset_id"`
}

package storage

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/domain"
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
	DueDate       *time.Time        `json:"due_date,omitempty"`
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

// Implements TaskDomain interface in domain
// But contains more detailed information of material, area and crop.
type TaskDomainDetailedCrop struct {
	Material *TaskDomainCropMaterial `json:"material"`
	Area     *TaskDomainCropArea     `json:"area"`
	Crop     *TaskDomainCropBatch    `json:"crop"`
}

type TaskDomainCropArea struct {
	AreaID   *uuid.UUID `json:"area_id"`
	AreaName string     `json:"area_name"`
}

type TaskDomainCropBatch struct {
	CropID      *uuid.UUID `json:"crop_id"`
	CropBatchID string     `json:"crop_batch_id"`
}

type TaskDomainCropMaterial struct {
	MaterialID           *uuid.UUID `json:"material_id"`
	MaterialName         string     `json:"material_name"`
	MaterialType         string     `json:"material_type"`
	MaterialDetailedType string     `json:"material_detailed_type"`
}

func (TaskDomainDetailedCrop) Code() string {
	return domain.TaskDomainCropCode
}

type TaskDomainDetailedArea struct {
	MaterialID           *uuid.UUID `json:"material_id"`
	MaterialName         string     `json:"material_name"`
	MaterialType         string     `json:"material_type"`
	MaterialDetailedType string     `json:"material_detailed_type"`
}

func (TaskDomainDetailedArea) Code() string {
	return domain.TaskDomainCropCode
}

type TaskDomainDetailedReservoir struct {
	MaterialID           *uuid.UUID `json:"material_id"`
	MaterialName         string     `json:"material_name"`
	MaterialType         string     `json:"material_type"`
	MaterialDetailedType string     `json:"material_detailed_type"`
}

func (TaskDomainDetailedReservoir) Code() string {
	return domain.TaskDomainCropCode
}

package domain

import (
	"github.com/gofrs/uuid"
)

const (
	TaskDomainAreaCode      = "AREA"
	TaskDomainCropCode      = "CROP"
	TaskDomainFinanceCode   = "FINANCE"
	TaskDomainGeneralCode   = "GENERAL"
	TaskDomainInventoryCode = "INVENTORY"
	TaskDomainReservoirCode = "RESERVOIR"
)

type TaskDomain interface {
	Code() string
}

// AREA
type TaskDomainArea struct {
	MaterialID *uuid.UUID `json:"material_id"`
}

func (TaskDomainArea) Code() string {
	return TaskDomainAreaCode
}

// CROP
type TaskDomainCrop struct {
	MaterialID *uuid.UUID `json:"material_id"`
	AreaID     *uuid.UUID `json:"area_id"`
}

func (TaskDomainCrop) Code() string {
	return TaskDomainCropCode
}

// FINANCE
type TaskDomainFinance struct{}

func (TaskDomainFinance) Code() string {
	return TaskDomainFinanceCode
}

// GENERAL
type TaskDomainGeneral struct{}

func (TaskDomainGeneral) Code() string {
	return TaskDomainGeneralCode
}

// INVENTORY
type TaskDomainInventory struct{}

func (TaskDomainInventory) Code() string {
	return TaskDomainInventoryCode
}

// RESERVOIR
type TaskDomainReservoir struct {
	MaterialID *uuid.UUID `json:"material_id"`
}

func (TaskDomainReservoir) Code() string {
	return TaskDomainReservoirCode
}

// CreateTaskDomainArea
func CreateTaskDomainArea(taskService TaskService, category string, materialID *uuid.UUID) (TaskDomainArea, error) {
	err := validateTaskCategory(category)
	if err != nil {
		return TaskDomainArea{}, err
	}

	if materialID != nil {
		err := validateAssetID(taskService, materialID, TaskDomainInventoryCode)
		if err != nil {
			return TaskDomainArea{}, err
		}
	}

	return TaskDomainArea{
		MaterialID: materialID,
	}, nil
}

// CreateTaskDomainCrop
func CreateTaskDomainCrop(ts TaskService, category string, materialID, areaID *uuid.UUID) (TaskDomainCrop, error) {
	err := validateTaskCategory(category)
	if err != nil {
		return TaskDomainCrop{}, err
	}

	if materialID != nil {
		err := validateAssetID(ts, materialID, TaskDomainInventoryCode)
		if err != nil {
			return TaskDomainCrop{}, err
		}
	}

	if areaID != nil {
		err := validateAssetID(ts, areaID, TaskDomainAreaCode)
		if err != nil {
			return TaskDomainCrop{}, err
		}
	}

	return TaskDomainCrop{
		MaterialID: materialID,
		AreaID:     areaID,
	}, nil
}

// CreateTaskDomainFinance
func CreateTaskDomainFinance() (TaskDomainFinance, error) {
	return TaskDomainFinance{}, nil
}

// CreateTaskDomainGeneral
func CreateTaskDomainGeneral() (TaskDomainGeneral, error) {
	return TaskDomainGeneral{}, nil
}

// CreateTaskDomainInventory
func CreateTaskDomainInventory() (TaskDomainInventory, error) {
	return TaskDomainInventory{}, nil
}

// CreateTaskDomainReservoir
func CreateTaskDomainReservoir(ts TaskService, category string, materialID *uuid.UUID) (TaskDomainReservoir, error) {
	err := validateTaskCategory(category)
	if err != nil {
		return TaskDomainReservoir{}, err
	}

	if materialID != nil {
		err := validateAssetID(ts, materialID, TaskDomainInventoryCode)
		if err != nil {
			return TaskDomainReservoir{}, err
		}
	}

	return TaskDomainReservoir{
		MaterialID: materialID,
	}, nil
}

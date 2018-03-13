package domain

import (
	uuid "github.com/satori/go.uuid"
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

func (d TaskDomainArea) Code() string {
	return TaskDomainAreaCode
}

// CROP
type TaskDomainCrop struct {
	MaterialID *uuid.UUID `json:"material_id"`
	AreaID     *uuid.UUID `json:"area_id"`
}

func (d TaskDomainCrop) Code() string {
	return TaskDomainCropCode
}

// FINANCE
type TaskDomainFinance struct {
}

func (d TaskDomainFinance) Code() string {
	return TaskDomainFinanceCode
}

// GENERAL
type TaskDomainGeneral struct {
}

func (d TaskDomainGeneral) Code() string {
	return TaskDomainGeneralCode
}

// INVENTORY
type TaskDomainInventory struct {
}

func (d TaskDomainInventory) Code() string {
	return TaskDomainInventoryCode
}

// RESERVOIR
type TaskDomainReservoir struct {
	MaterialID *uuid.UUID `json:"material_id"`
}

func (d TaskDomainReservoir) Code() string {
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
func CreateTaskDomainCrop(taskService TaskService, category string, materialID *uuid.UUID, areaID *uuid.UUID) (TaskDomainCrop, error) {

	err := validateTaskCategory(category)
	if err != nil {
		return TaskDomainCrop{}, err
	}

	if materialID != nil {
		err := validateAssetID(taskService, materialID, TaskDomainInventoryCode)
		if err != nil {
			return TaskDomainCrop{}, err
		}
	}

	if areaID != nil {
		err := validateAssetID(taskService, areaID, TaskDomainAreaCode)
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
func CreateTaskDomainReservoir(taskService TaskService, category string, materialID *uuid.UUID) (TaskDomainReservoir, error) {

	err := validateTaskCategory(category)
	if err != nil {
		return TaskDomainReservoir{}, err
	}

	if materialID != nil {
		err := validateAssetID(taskService, materialID, TaskDomainInventoryCode)
		if err != nil {
			return TaskDomainReservoir{}, err
		}
	}

	return TaskDomainReservoir{
		MaterialID: materialID,
	}, nil
}

// validateAssetID
func validateDomainAssetID(taskService TaskService, assetid *uuid.UUID, taskdomain string) error {
	if taskdomain == "" {
		return TaskError{TaskErrorDomainEmptyCode}
	}
	//Find asset in repository
	// if not found return error

	switch taskdomain {
	case TaskDomainAreaCode:

		serviceResult := taskService.FindAreaByID(*assetid)

		if serviceResult.Error != nil {
			return serviceResult.Error
		}
	case TaskDomainCropCode:

		serviceResult := taskService.FindCropByID(*assetid)

		if serviceResult.Error != nil {
			return serviceResult.Error
		}
	case TaskDomainInventoryCode:

		serviceResult := taskService.FindMaterialByID(*assetid)

		if serviceResult.Error != nil {
			return serviceResult.Error
		}
	case TaskDomainReservoirCode:

		serviceResult := taskService.FindReservoirByID(*assetid)

		if serviceResult.Error != nil {
			return serviceResult.Error
		}

	default:
		return TaskError{TaskErrorInvalidDomainCode}
	}
	return nil
}

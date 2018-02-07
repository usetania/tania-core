package service

import (
	domain "github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

// TaskService handles task behaviours that needs external interaction to be worked

type TaskServiceInMemory struct {
	CropQuery     query.CropQuery
	AreaQuery     query.AreaQuery
	MaterialQuery query.MaterialQuery
}

func (s TaskServiceInMemory) FindAreaByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.AreaQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	area, ok := result.Result.(query.TaskAreaQueryResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if area == (query.TaskAreaQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: area,
	}
}

func (s TaskServiceInMemory) FindCropByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.CropQuery.FindCropByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	crop, ok := result.Result.(query.TaskCropQueryResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if crop == (query.TaskCropQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: crop,
	}
}

func (s TaskServiceInMemory) FindMaterialByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.MaterialQuery.FindMaterialByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	material, ok := result.Result.(query.TaskMaterialQueryResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if material == (query.TaskMaterialQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: material,
	}
}

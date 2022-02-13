package service

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/tasks/domain"
	"github.com/usetania/tania-core/src/tasks/query"
)

// TaskService handles task behaviours that needs external interaction to be worked

type TaskServiceSqlite struct {
	CropQuery      query.Crop
	AreaQuery      query.Area
	MaterialQuery  query.Material
	ReservoirQuery query.Reservoir
}

func (s TaskServiceSqlite) FindAreaByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.AreaQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	area, ok := result.Result.(query.TaskAreaResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if area == (query.TaskAreaResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: area,
	}
}

func (s TaskServiceSqlite) FindCropByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.CropQuery.FindCropByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	crop, ok := result.Result.(query.TaskCropResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if crop == (query.TaskCropResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: crop,
	}
}

func (s TaskServiceSqlite) FindMaterialByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.MaterialQuery.FindMaterialByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	material, ok := result.Result.(query.TaskMaterialResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if material == (query.TaskMaterialResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: material,
	}
}

func (s TaskServiceSqlite) FindReservoirByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.ReservoirQuery.FindReservoirByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	reservoir, ok := result.Result.(query.TaskReservoirResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if reservoir == (query.TaskReservoirResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: reservoir,
	}
}

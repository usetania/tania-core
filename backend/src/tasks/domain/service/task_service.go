package service

import (
	domain "github.com/Tanibox/tania-core/src/tasks/domain"
	"github.com/Tanibox/tania-core/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

// TaskService handles task behaviours that needs external interaction to be worked

type TaskServiceSqlLite struct {
	CropQuery      query.CropQuery
	AreaQuery      query.AreaQuery
	MaterialQuery  query.MaterialQuery
	ReservoirQuery query.ReservoirQuery
}

func (s TaskServiceSqlLite) FindAreaByID(uid uuid.UUID) domain.ServiceResult {
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

func (s TaskServiceSqlLite) FindCropByID(uid uuid.UUID) domain.ServiceResult {
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

func (s TaskServiceSqlLite) FindMaterialByID(uid uuid.UUID) domain.ServiceResult {
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

func (s TaskServiceSqlLite) FindReservoirByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.ReservoirQuery.FindReservoirByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	reservoir, ok := result.Result.(query.TaskReservoirQueryResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	if reservoir == (query.TaskReservoirQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.TaskError{Code: domain.TaskErrorInvalidAssetIDCode},
		}
	}

	return domain.ServiceResult{
		Result: reservoir,
	}
}

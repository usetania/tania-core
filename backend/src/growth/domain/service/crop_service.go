package service

import (
	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/growth/domain"
	"github.com/usetania/tania-core/src/growth/query"
	"github.com/usetania/tania-core/src/growth/storage"
)

type CropServiceInMemory struct {
	MaterialReadQuery query.MaterialReadQuery
	CropReadQuery     query.CropReadQuery
	AreaReadQuery     query.AreaReadQuery
}

func (s CropServiceInMemory) FindMaterialByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.MaterialReadQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	inv, ok := result.Result.(query.CropMaterialQueryResult)

	if !ok {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropMaterialErrorInvalidMaterial},
		}
	}

	if inv == (query.CropMaterialQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropMaterialErrorNotFound},
		}
	}

	return domain.ServiceResult{
		Result: inv,
	}
}

func (s CropServiceInMemory) FindByBatchID(batchID string) domain.ServiceResult {
	resultQuery := <-s.CropReadQuery.FindByBatchID(batchID)

	if resultQuery.Error != nil {
		return domain.ServiceResult{
			Error: resultQuery.Error,
		}
	}

	cropFound, ok := resultQuery.Result.(storage.CropRead)
	if !ok {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropErrorInvalidBatchID},
		}
	}

	if cropFound.UID != (uuid.UUID{}) {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropErrorBatchIDAlreadyCreated},
		}
	}

	return domain.ServiceResult{
		Result: cropFound,
	}
}

func (s CropServiceInMemory) FindAreaByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.AreaReadQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	area, ok := result.Result.(query.CropAreaQueryResult)
	if !ok {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropMoveToAreaErrorInvalidSourceArea},
		}
	}

	if area == (query.CropAreaQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropMoveToAreaErrorSourceAreaNotFound},
		}
	}

	return domain.ServiceResult{
		Result: area,
	}
}

package service

import (
	"github.com/Tanibox/tania-server/src/growth/domain"
	"github.com/Tanibox/tania-server/src/growth/query"
	uuid "github.com/satori/go.uuid"
)

type CropServiceInMemory struct {
	InventoryMaterialQuery query.InventoryMaterialQuery
	CropQuery              query.CropQuery
	AreaQuery              query.AreaQuery
}

func (s CropServiceInMemory) FindInventoryMaterialByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.InventoryMaterialQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{
			Error: result.Error,
		}
	}

	inv, ok := result.Result.(query.CropInventoryQueryResult)

	if !ok {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropInventoryErrorInvalidInventory},
		}
	}

	if inv == (query.CropInventoryQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.CropError{Code: domain.CropInventoryErrorNotFound},
		}
	}

	return domain.ServiceResult{
		Result: inv,
	}
}

func (s CropServiceInMemory) FindByBatchID(batchID string) domain.ServiceResult {
	resultQuery := <-s.CropQuery.FindByBatchID(batchID)

	if resultQuery.Error != nil {
		return domain.ServiceResult{
			Error: resultQuery.Error,
		}
	}

	cropFound, ok := resultQuery.Result.(domain.Crop)
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
	result := <-s.AreaQuery.FindByID(uid)

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

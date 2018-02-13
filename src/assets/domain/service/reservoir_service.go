package service

import (
	"github.com/Tanibox/tania-server/src/assets/domain"
	"github.com/Tanibox/tania-server/src/assets/query"
	uuid "github.com/satori/go.uuid"
)

type Farm struct {
	UID uuid.UUID
}

type ReservoirServiceInMemory struct {
	FarmReadQuery query.FarmReadQuery
}

func (s ReservoirServiceInMemory) FindFarmByID(uid uuid.UUID) domain.ServiceResult {
	result := <-s.FarmReadQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ServiceResult{Error: result.Error}
	}

	farm, ok := result.Result.(query.FarmReadQueryResult)

	if !ok {
		return domain.ServiceResult{
			Error: domain.ReservoirError{Code: domain.ReservoirErrorFarmNotFound},
		}
	}

	if farm == (query.FarmReadQueryResult{}) {
		return domain.ServiceResult{
			Error: domain.ReservoirError{Code: domain.ReservoirErrorFarmNotFound},
		}
	}

	return domain.ServiceResult{
		Result: farm,
	}
}

package service

import (
	"github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/Tanibox/tania-core/src/assets/query"
	"github.com/Tanibox/tania-core/src/assets/storage"
	"github.com/gofrs/uuid"
)

type ReservoirServiceInMemory struct {
	FarmReadQuery query.FarmReadQuery
}

func (s ReservoirServiceInMemory) FindFarmByID(uid uuid.UUID) (domain.ReservoirFarmServiceResult, error) {
	result := <-s.FarmReadQuery.FindByID(uid)

	if result.Error != nil {
		return domain.ReservoirFarmServiceResult{}, result.Error
	}

	farm, ok := result.Result.(storage.FarmRead)

	if !ok {
		return domain.ReservoirFarmServiceResult{}, domain.ReservoirError{Code: domain.ReservoirErrorFarmNotFound}
	}

	if farm == (storage.FarmRead{}) {
		return domain.ReservoirFarmServiceResult{}, domain.ReservoirError{Code: domain.ReservoirErrorFarmNotFound}
	}

	return domain.ReservoirFarmServiceResult{
		UID:  farm.UID,
		Name: farm.Name,
	}, nil
}

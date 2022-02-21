package service

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/query"
	"github.com/usetania/tania-core/src/assets/storage"
)

type AreaServiceInMemory struct {
	FarmReadQuery      query.FarmRead
	ReservoirReadQuery query.ReservoirRead
	CropReadQuery      query.CropRead
}

func (s AreaServiceInMemory) FindFarmByID(uid uuid.UUID) (domain.AreaFarmServiceResult, error) {
	result := <-s.FarmReadQuery.FindByID(uid)

	if result.Error != nil {
		return domain.AreaFarmServiceResult{}, result.Error
	}

	farm, ok := result.Result.(storage.FarmRead)

	if !ok {
		return domain.AreaFarmServiceResult{}, domain.AreaError{Code: domain.AreaErrorFarmNotFound}
	}

	if farm == (storage.FarmRead{}) {
		return domain.AreaFarmServiceResult{}, domain.AreaError{Code: domain.AreaErrorFarmNotFound}
	}

	return domain.AreaFarmServiceResult{
		UID:  farm.UID,
		Name: farm.Name,
	}, nil
}

func (s AreaServiceInMemory) FindReservoirByID(reservoirUID uuid.UUID) (domain.AreaReservoirServiceResult, error) {
	result := <-s.ReservoirReadQuery.FindByID(reservoirUID)

	if result.Error != nil {
		return domain.AreaReservoirServiceResult{}, result.Error
	}

	res, ok := result.Result.(storage.ReservoirRead)

	if !ok {
		return domain.AreaReservoirServiceResult{}, domain.AreaError{Code: domain.AreaErrorReservoirNotFound}
	}

	if res.UID == (uuid.UUID{}) {
		return domain.AreaReservoirServiceResult{}, domain.AreaError{Code: domain.AreaErrorReservoirNotFound}
	}

	return domain.AreaReservoirServiceResult{
		UID:  res.UID,
		Name: res.Name,
	}, nil
}

func (s AreaServiceInMemory) CountCropsByAreaID(areaUID uuid.UUID) (int, error) {
	result := <-s.CropReadQuery.CountCropsByArea(areaUID)
	if result.Error != nil {
		return 0, result.Error
	}

	totals, ok := result.Result.(query.CountAreaCropResult)
	if !ok {
		return 0, errors.New("internal server error")
	}

	return totals.TotalCropBatch, nil
}

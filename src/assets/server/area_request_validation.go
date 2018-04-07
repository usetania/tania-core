package server

import (
	"strconv"

	"github.com/Tanibox/tania-core/src/assets/domain"
	"github.com/Tanibox/tania-core/src/assets/storage"
	uuid "github.com/satori/go.uuid"
)

func (rv *RequestValidation) ValidateReservoir(s FarmServer, reservoirUID uuid.UUID) (storage.ReservoirRead, error) {
	result := <-s.ReservoirReadQuery.FindByID(reservoirUID)
	reservoir, _ := result.Result.(storage.ReservoirRead)

	if reservoir.UID == (uuid.UUID{}) {
		return reservoir, NewRequestValidationError(NOT_FOUND, "reservoir_id")
	}

	return reservoir, nil
}

func (rv *RequestValidation) ValidateFarm(s FarmServer, farmUID uuid.UUID) (storage.FarmRead, error) {
	result := <-s.FarmReadQuery.FindByID(farmUID)
	farm, _ := result.Result.(storage.FarmRead)

	if farm.UID == (uuid.UUID{}) {
		return farm, NewRequestValidationError(NOT_FOUND, "farm_id")
	}

	return farm, nil
}

func (rv *RequestValidation) ValidateAreaSize(size string, sizeUnit string) (domain.AreaSize, error) {
	if size == "" {
		return domain.AreaSize{}, NewRequestValidationError(REQUIRED, "size")
	}

	if sizeUnit == "" {
		return domain.AreaSize{}, NewRequestValidationError(REQUIRED, "size_unit")
	}

	sizeFloat, err := strconv.ParseFloat(size, 32)
	if err != nil {
		return domain.AreaSize{}, err
	}

	unit := domain.GetAreaUnit(sizeUnit)
	if unit == (domain.AreaUnit{}) {
		return domain.AreaSize{}, NewRequestValidationError(INVALID_OPTION, "size_unit")
	}

	return domain.AreaSize{
		Value: float32(sizeFloat),
		Unit:  unit,
	}, nil
}

func (rv *RequestValidation) ValidateAreaLocation(location string) (string, error) {
	if location == "" {
		return "", NewRequestValidationError(REQUIRED, "location")
	}

	areaLocation := domain.GetAreaLocation(location)
	if areaLocation == (domain.AreaLocation{}) {
		return "", NewRequestValidationError(INVALID_OPTION, "location")
	}

	return areaLocation.Code, nil
}

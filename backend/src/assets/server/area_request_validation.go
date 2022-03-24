package server

import (
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/assets/storage"
)

func (*RequestValidation) ValidateReservoir(s FarmServer, reservoirUID uuid.UUID) (storage.ReservoirRead, error) {
	result := <-s.ReservoirReadQuery.FindByID(reservoirUID)
	reservoir, _ := result.Result.(storage.ReservoirRead)

	if reservoir.UID == (uuid.UUID{}) {
		return reservoir, NewRequestValidationError(NotFound, "reservoir_id")
	}

	return reservoir, nil
}

func (*RequestValidation) ValidateFarm(s FarmServer, farmUID uuid.UUID) (storage.FarmRead, error) {
	result := <-s.FarmReadQuery.FindByID(farmUID)
	farm, _ := result.Result.(storage.FarmRead)

	if farm.UID == (uuid.UUID{}) {
		return farm, NewRequestValidationError(NotFound, "farm_id")
	}

	return farm, nil
}

func (*RequestValidation) ValidateAreaSize(size, sizeUnit string) (domain.AreaSize, error) {
	if size == "" {
		return domain.AreaSize{}, NewRequestValidationError(Required, "size")
	}

	if sizeUnit == "" {
		return domain.AreaSize{}, NewRequestValidationError(Required, "size_unit")
	}

	sizeFloat, err := strconv.ParseFloat(size, 32)
	if err != nil {
		return domain.AreaSize{}, err
	}

	unit := domain.GetAreaUnit(sizeUnit)
	if unit == (domain.AreaUnit{}) {
		return domain.AreaSize{}, NewRequestValidationError(InvalidOption, "size_unit")
	}

	return domain.AreaSize{
		Value: float32(sizeFloat),
		Unit:  unit,
	}, nil
}

func (*RequestValidation) ValidateAreaLocation(location string) (string, error) {
	if location == "" {
		return "", NewRequestValidationError(Required, "location")
	}

	areaLocation := domain.GetAreaLocation(location)
	if areaLocation == (domain.AreaLocation{}) {
		return "", NewRequestValidationError(InvalidOption, "location")
	}

	return areaLocation.Code, nil
}

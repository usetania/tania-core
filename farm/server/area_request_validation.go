package server

import (
	"strconv"

	"github.com/Tanibox/tania-server/farm/entity"
)

func (rv *RequestValidation) ValidateReservoir(s FarmServer, farmId string) (entity.Reservoir, error) {
	result := <-s.ReservoirRepo.FindByID(farmId)
	reservoir, _ := result.Result.(entity.Reservoir)

	if reservoir.UID == "" {
		return reservoir, NewRequestValidationError(NOT_FOUND, "reservoir_id")
	}

	return reservoir, nil
}

func (rv *RequestValidation) ValidateAreaSize(size string, sizeUnit string) (entity.AreaUnit, error) {
	if size == "" {
		return nil, NewRequestValidationError(REQUIRED, "size")
	}

	if sizeUnit == "" {
		return nil, NewRequestValidationError(REQUIRED, "size_unit")
	}

	sizeFloat, err := strconv.ParseFloat(size, 32)
	if err != nil {
		return nil, err
	}

	switch sizeUnit {
	case "m2":
		return entity.SquareMeter{Value: float32(sizeFloat)}, nil
	case "hectare":
		return entity.Hectare{Value: float32(sizeFloat)}, nil
	default:
		return nil, NewRequestValidationError(INVALID_OPTION, "size_unit")
	}
}

func (rv *RequestValidation) ValidateAreaLocation(location string) (string, error) {
	if location == "" {
		return "", NewRequestValidationError(REQUIRED, "location")
	}

	areaLocation, err := entity.FindAreaLocationByCode(location)
	if err != nil {
		return "", NewRequestValidationError(INVALID_OPTION, "location")
	}

	return areaLocation.Code, nil
}

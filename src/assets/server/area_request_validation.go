package server

import (
	"strconv"

	"github.com/Tanibox/tania-server/src/assets/domain"
	uuid "github.com/satori/go.uuid"
)

func (rv *RequestValidation) ValidateReservoir(s FarmServer, reservoirId string) (domain.Reservoir, error) {
	result := <-s.ReservoirRepo.FindByID(reservoirId)
	reservoir, _ := result.Result.(domain.Reservoir)

	if reservoir.UID == (uuid.UUID{}) {
		return reservoir, NewRequestValidationError(NOT_FOUND, "reservoir_id")
	}

	return reservoir, nil
}

func (rv *RequestValidation) ValidateAreaSize(size string, sizeUnit string) (domain.AreaUnit, error) {
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
		return domain.SquareMeter{Value: float32(sizeFloat)}, nil
	case "hectare":
		return domain.Hectare{Value: float32(sizeFloat)}, nil
	default:
		return nil, NewRequestValidationError(INVALID_OPTION, "size_unit")
	}
}

func (rv *RequestValidation) ValidateAreaLocation(location string) (string, error) {
	if location == "" {
		return "", NewRequestValidationError(REQUIRED, "location")
	}

	areaLocation, err := domain.FindAreaLocationByCode(location)
	if err != nil {
		return "", NewRequestValidationError(INVALID_OPTION, "location")
	}

	return areaLocation.Code, nil
}

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

func (rv *RequestValidation) ValidateAreaSize(size string) (float32, error) {
	if size == "" {
		return 0, NewRequestValidationError(REQUIRED, "size")
	}

	sizeFloat, err := strconv.ParseFloat(size, 32)

	return float32(sizeFloat), err
}

func (rv *RequestValidation) ValidateAreaSizeUnit(sizeUnit string) (string, error) {
	if sizeUnit == "" {
		return "", NewRequestValidationError(REQUIRED, "size_unit")
	}

	// TODO: Need to be refactored
	if sizeUnit != "m2" && sizeUnit != "hectar" {
		return "", NewRequestValidationError(INVALID_OPTION, "size_unit")
	}

	return sizeUnit, nil
}

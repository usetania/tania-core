package server

import (
	"strconv"

	"github.com/Tanibox/tania-server/src/helper/validationhelper"

	"github.com/Tanibox/tania-server/src/assets/domain"
)

func (rv *RequestValidation) ValidateReservoirName(name string) (string, error) {
	if name == "" {
		return "", NewRequestValidationError(REQUIRED, "name")
	}
	if !validationhelper.IsAlphanumeric(name) {
		return "", NewRequestValidationError(ALPHANUMERIC, "name")
	}

	return name, nil
}

func (rv *RequestValidation) ValidateCapacity(waterSourceType, capacity string) (float32, error) {
	if waterSourceType == "tap" {
		return 0, nil
	}

	if capacity == "" {
		return 0, NewRequestValidationError(REQUIRED, "capacity")
	}

	if !validationhelper.IsFloat(capacity) {
		return 0, NewRequestValidationError(FLOAT, "capacity")
	}

	c, err := strconv.ParseFloat(capacity, 32)
	if err != nil {
		return 0, NewRequestValidationError(PARSE_FAILED, "capacity")
	}

	return float32(c), nil
}

func (rv *RequestValidation) ValidateType(t string) (string, error) {
	if t == "" {
		return "", NewRequestValidationError(REQUIRED, "type")
	}

	if !validationhelper.IsAlpha(t) {
		return "", NewRequestValidationError(ALPHA, "type")
	}

	if t != "bucket" && t != "tap" {
		return "", NewRequestValidationError(INVALID_OPTION, "type")
	}

	return t, nil
}

func (rv *RequestValidation) ValidateFarm(s FarmServer, farmId string) (domain.Farm, error) {
	result := <-s.FarmRepo.FindByID(farmId)
	farm, _ := result.Result.(domain.Farm)

	if farm.UID.String() == "" {
		return domain.Farm{}, NewRequestValidationError(NOT_FOUND, "farm_id")
	}

	return farm, nil
}

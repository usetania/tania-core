package server

import (
	"fmt"
	"strconv"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/Tanibox/tania-server/helper/validationhelper"
)

func (rve RequestValidationError) Error() string {
	return fmt.Sprintf(
		"Field Name: %s, Error Code: %s, Error Message: %s",
		rve.FieldName,
		rve.ErrorCode,
		rve.ErrorMessage,
	)
}

// RequestValidation sanitizes request inputs and convert the input to its correct data type.
// This is mostly used to prevent issues like invalid data type or potential SQL Injection.
// So we can focus on processing data without converting data type after this sanitizing.
// This validation doesn't aim to validate business process.
// The business process validation will be handled in each entity's behaviour.
type RequestValidation struct {
}

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

func (rv *RequestValidation) ValidateFarm(farmId string) (entity.Farm, error) {
	return entity.Farm{}, nil
}

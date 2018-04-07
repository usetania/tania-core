package server

import (
	"strconv"

	"github.com/Tanibox/tania-core/src/helper/validationhelper"

	"github.com/Tanibox/tania-core/src/assets/domain"
)

func (rv *RequestValidation) ValidateReservoirName(name string) (string, error) {
	if name == "" {
		return "", NewRequestValidationError(REQUIRED, "name")
	}
	if !validationhelper.IsAlphanumSpaceHyphenUnderscore(name) {
		return "", NewRequestValidationError(ALPHANUMERIC, "name")
	}

	return name, nil
}

func (rv *RequestValidation) ValidateCapacity(waterSourceType, capacity string) (float32, error) {
	if waterSourceType == domain.TapType {
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

	if t != domain.BucketType && t != domain.TapType {
		return "", NewRequestValidationError(INVALID_OPTION, "type")
	}

	return t, nil
}

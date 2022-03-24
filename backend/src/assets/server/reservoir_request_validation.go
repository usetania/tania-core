package server

import (
	"strconv"

	"github.com/usetania/tania-core/src/assets/domain"
	"github.com/usetania/tania-core/src/helper/validationhelper"
)

func (*RequestValidation) ValidateReservoirName(name string) (string, error) {
	if name == "" {
		return "", NewRequestValidationError(Required, "name")
	}

	if !validationhelper.IsAlphanumSpaceHyphenUnderscore(name) {
		return "", NewRequestValidationError(Alphanumeric, "name")
	}

	return name, nil
}

func (*RequestValidation) ValidateCapacity(waterSourceType, capacity string) (float32, error) {
	if waterSourceType == domain.TapType {
		return 0, nil
	}

	if capacity == "" {
		return 0, NewRequestValidationError(Required, "capacity")
	}

	if !validationhelper.IsFloat(capacity) {
		return 0, NewRequestValidationError(Float, "capacity")
	}

	c, err := strconv.ParseFloat(capacity, 32)
	if err != nil {
		return 0, NewRequestValidationError(ParseFailed, "capacity")
	}

	return float32(c), nil
}

func (*RequestValidation) ValidateType(t string) (string, error) {
	if t == "" {
		return "", NewRequestValidationError(Required, "type")
	}

	if !validationhelper.IsAlpha(t) {
		return "", NewRequestValidationError(Alpha, "type")
	}

	if t != domain.BucketType && t != domain.TapType {
		return "", NewRequestValidationError(InvalidOption, "type")
	}

	return t, nil
}

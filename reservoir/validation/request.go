package validation

import (
	"fmt"
	"strconv"

	"github.com/Tanibox/tania-server/helper/validationhelper"
)

const (
	REQUIRED       = "REQUIRED"
	ALPHANUMERIC   = "ALPHANUMERIC"
	ALPHA          = "ALPHA"
	NUMERIC        = "NUMERIC"
	FLOAT          = "FLOAT"
	PARSE_FAILED   = "PARSE_FAILED"
	INVALID_OPTION = "INVALID_OPTION"
)

type RequestValidationError struct {
	FieldName    string `json:"field_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func (rve RequestValidationError) Error() string {
	return fmt.Sprintf(
		"Field Name: %s, Error Code: %s, Error Message: %s",
		rve.FieldName,
		rve.ErrorCode,
		rve.ErrorMessage,
	)
}

func Message(errorCode string) string {
	switch errorCode {
	case REQUIRED:
		return "This field is required"
	case ALPHANUMERIC:
		return "Alphanumeric only"
	case ALPHA:
		return "Alphabet only"
	case NUMERIC:
		return "Number only"
	case FLOAT:
		return "Float only"
	case PARSE_FAILED:
		return "Parsing failed. Make sure the input is correct."
	case INVALID_OPTION:
		return "This value is not available in options. Please give the correct options."
	default:
		return "Internal server error"
	}
}

func NewRequestValidationError(errorCode, fieldName string) RequestValidationError {
	return RequestValidationError{
		FieldName:    fieldName,
		ErrorCode:    errorCode,
		ErrorMessage: Message(errorCode),
	}
}

// RequestValidation sanitizes request inputs and convert the input to its correct data type.
// This is mostly used to prevent issues like invalid data type or potential SQL Injection.
// So we can focus on processing data without converting data type after this sanitizing.
// This validation doesn't aim to validate business process.
// The business process validation will be handled in each entity's behaviour.
type RequestValidation struct {
}

func (rv *RequestValidation) ValidateName(name string) (string, error) {
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

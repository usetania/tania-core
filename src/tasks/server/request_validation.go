package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/src/tasks/domain"
	"github.com/labstack/echo"
)

const (
	REQUIRED       = "REQUIRED"
	ALPHANUMERIC   = "ALPHANUMERIC"
	ALPHA          = "ALPHA"
	NUMERIC        = "NUMERIC"
	FLOAT          = "FLOAT"
	PARSE_FAILED   = "PARSE_FAILED"
	INVALID_OPTION = "INVALID_OPTION"
	NOT_FOUND      = "NOT_FOUND"
)

// RequestValidation sanitizes request inputs and convert the input to its correct data type.
// This is mostly used to prevent issues like invalid data type or potential SQL Injection.
// So we can focus on processing data without converting data type after this sanitizing.
// This validation doesn't aim to validate business process.
// The business process validation will be handled in each entity's behaviour.
type RequestValidation struct {
}

// RequestValidationError contains fields used for JSON error response
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

// Message translates error code to meaningful message
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
	case NOT_FOUND:
		return "Data not found."
	default:
		return "Internal server error"
	}
}

// NewRequestValidationError initializes new RequestValidation struct
func NewRequestValidationError(errorCode, fieldName string) RequestValidationError {
	return RequestValidationError{
		FieldName:    fieldName,
		ErrorCode:    errorCode,
		ErrorMessage: Message(errorCode),
	}
}

// Error wraps errors from application layer and domain layer
// to some format in JSON for response
func Error(c echo.Context, err error) error {
	errorResponse := map[string]string{
		"field_name":    "",
		"error_code":    "",
		"error_message": "",
	}

	if re, ok := err.(domain.ReservoirError); ok {
		errorResponse["error_code"] = strconv.Itoa(re.Code)
		errorResponse["error_message"] = re.Error()

		return c.JSON(http.StatusBadRequest, errorResponse)
	} else if re, ok := err.(domain.FarmError); ok {
		errorResponse["error_code"] = strconv.Itoa(re.Code)
		errorResponse["error_message"] = re.Error()

		return c.JSON(http.StatusBadRequest, errorResponse)
	} else if re, ok := err.(domain.AreaError); ok {
		errorResponse["error_code"] = strconv.Itoa(re.Code)
		errorResponse["error_message"] = re.Error()

		return c.JSON(http.StatusBadRequest, errorResponse)
	} else if rve, ok := err.(RequestValidationError); ok {
		errorResponse["field_name"] = rve.FieldName
		errorResponse["error_code"] = rve.ErrorCode
		errorResponse["error_message"] = rve.ErrorMessage

		return c.JSON(http.StatusBadRequest, rve)
	}

	errorResponse["error_message"] = err.Error()
	return c.JSON(http.StatusInternalServerError, errorResponse)
}

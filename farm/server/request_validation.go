package server

import (
	"net/http"
	"strconv"

	"github.com/Tanibox/tania-server/farm/entity"
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
)

// RequestValidationError contains fields used for JSON error response
type RequestValidationError struct {
	FieldName    string `json:"field_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
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
	if re, ok := err.(entity.ReservoirError); ok {
		errMap := map[string]string{
			"field_name":    "",
			"error_code":    strconv.Itoa(re.Code),
			"error_message": re.Error(),
		}

		return c.JSON(http.StatusBadRequest, errMap)
	} else if re, ok := err.(entity.FarmError); ok {
		errMap := map[string]string{
			"field_name":    "",
			"error_code":    strconv.Itoa(re.Code),
			"error_message": re.Error(),
		}

		return c.JSON(http.StatusBadRequest, errMap)
	} else if rve, ok := err.(RequestValidationError); ok {
		return c.JSON(http.StatusBadRequest, rve)
	}

	return c.JSON(http.StatusInternalServerError, err)
}

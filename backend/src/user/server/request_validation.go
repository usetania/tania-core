package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/usetania/tania-core/src/user/domain"
)

const (
	Required      = "REQUIRED"
	Alphanumeric  = "ALPHANUMERIC"
	Alpha         = "ALPHA"
	Numeric       = "NUMERIC"
	Float         = "FLOAT"
	ParseFailed   = "PARSE_FAILED"
	InvalidOption = "INVALID_OPTION"
	NotFound      = "NOT_FOUND"
	NorMatch      = "NOT_MATCH"
	Invalid       = "INVALID"
)

// RequestValidation sanitizes request inputs and convert the input to its correct data type.
// This is mostly used to prevent issues like invalid data type or potential SQL Injection.
// So we can focus on processing data without converting data type after this sanitizing.
// This validation doesn't aim to validate business process.
// The business process validation will be handled in each entity's behaviour.
type RequestValidation struct{}

// RequestValidationError contains fields used for JSON error response.
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

// Message translates error code to meaningful message.
func Message(errorCode string) string {
	switch errorCode {
	case Required:
		return "This field is required"
	case Alphanumeric:
		return "Alphanumeric only"
	case Alpha:
		return "Alphabet only"
	case Numeric:
		return "Number only"
	case Float:
		return "Float only"
	case ParseFailed:
		return "Parsing failed. Make sure the input is correct."
	case InvalidOption:
		return "This value is not available in options. Please give the correct options."
	case NotFound:
		return "Data not found."
	case NorMatch:
		return "Password didn't match with confirmation password"
	case Invalid:
		return "Invalid value"
	default:
		return "Internal server error"
	}
}

// NewRequestValidationError initializes new RequestValidation struct.
func NewRequestValidationError(errorCode, fieldName string) RequestValidationError {
	return RequestValidationError{
		FieldName:    fieldName,
		ErrorCode:    errorCode,
		ErrorMessage: Message(errorCode),
	}
}

// Error wraps errors from application layer and domain layer
// to some format in JSON for response.
func Error(c echo.Context, err error) error {
	errorResponse := map[string]string{
		"field_name":    "",
		"error_code":    "",
		"error_message": "",
	}

	file, line := getFileAndLineNumber()

	log.Printf(
		"user_uid: %v\nrequest_id: %v\nfile: %v\nline: %v\n",
		c.Get("USER_UID"),
		c.Response().Header().Get(echo.HeaderXRequestID),
		file,
		line,
	)

	errorResponse["error_message"] = err.Error()
	log.Printf("error_message: %v\n", err.Error())

	var ue domain.UserError
	if errors.As(err, &ue) {
		errorResponse["error_code"] = strconv.Itoa(ue.Code)

		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	var rve RequestValidationError
	if errors.As(err, &rve) {
		errorResponse["field_name"] = rve.FieldName
		errorResponse["error_code"] = rve.ErrorCode
		errorResponse["error_message"] = rve.ErrorMessage

		return c.JSON(http.StatusBadRequest, rve)
	}

	return c.JSON(http.StatusInternalServerError, errorResponse)
}

func getFileAndLineNumber() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	return file, line
}

package errors

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type APIError struct {
	StatusCode   int   `json:"status_code"`
	InternalCode int   `json:"internal_code"`
	Err          error `json:"-"`
}

func (e APIError) Error() string {
	return strconv.Itoa(e.InternalCode) + ": " + e.Err.Error()
}

func (e APIError) Unwrap() error {
	return e.Err
}

func (e APIError) Wrap(ogErr error) error {
	return WrapError(
		e.Error(),
		ogErr,
	)
}

func (e APIError) GetStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

func WrapError(msg string, ogerr error) *APIError {
	var apiError *APIError
	if errors.As(ogerr, &apiError) {
		return &APIError{
			Err:          errors.Wrap(apiError, msg),
			StatusCode:   apiError.StatusCode,
			InternalCode: apiError.InternalCode,
		}
	}
	return &APIError{
		Err:          errors.Wrap(ogerr, msg),
		StatusCode:   http.StatusInternalServerError,
		InternalCode: 999,
	}
}

var ErrorNoResponse = &APIError{
	StatusCode:   http.StatusInternalServerError,
	InternalCode: 1000,
	Err:          errors.New("no response"),
}

var ErrorNotFound = &APIError{
	StatusCode:   http.StatusNotFound,
	InternalCode: 1001,
	Err:          errors.New("not found"),
}

var ErrorBadBody = &APIError{
	StatusCode:   http.StatusBadRequest,
	InternalCode: 1002,
	Err:          errors.New("bad body"),
}

func NewJSONDecodeError(err error) error {
	if err == nil {
		return nil
	}
	var ute *json.UnmarshalTypeError
	var iue *json.InvalidUnmarshalError
	if errors.As(err, &ute) {
		return WrapError(ute.Error(), ErrorBadBody)
	} else if errors.As(err, &iue) {
		return WrapError(iue.Error(), ErrorBadBody)
	} else {
		return err
	}
}

func NewValidationError(err error) error {
	if err == nil {
		return nil
	}
	ves, ok := err.(validator.ValidationErrors)
	// if not a validation error, return the original error
	if !ok {
		return err
	}
	if len(ves) == 0 {
		return nil
	}
	return WrapError(ves.Error(), ErrorBadBody)
}

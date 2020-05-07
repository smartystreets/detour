package detour

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

func prepareInputModel(model interface{}, request *http.Request) (statusCode int, err error) {
	if err = Bind(request, model); err != nil {
		return statusCodeFromErrorOrDefault(err, http.StatusBadRequest)
	}

	sanitize(model)

	if err = validate(model); err != nil {
		return statusCodeFromErrorOrDefault(err, http.StatusUnprocessableEntity)
	}

	if err = serverError(model); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// Bind is exported for use in testing.
func Bind(request *http.Request, message interface{}) error {
	bindContext(message, request)

	err := bindJSON(request, message)
	if err != nil {
		return err
	}

	binder, isBinder := message.(Binder)
	if !isBinder {
		return nil
	}

	err = request.ParseForm()
	if err != nil {
		return err
	}

	err = binder.Bind(request)
	if err == nil {
		return nil
	}

	errs, isErrors := err.(Errors)
	if isErrors && len(errs) == 0 {
		return nil
	}

	diagnosticErrs, isDiagnosticErrors := err.(DiagnosticErrors)
	if isDiagnosticErrors && len(diagnosticErrs) == 0 {
		return nil
	}

	return err
}

func bindContext(message interface{}, request *http.Request) {
	if message == nil {
		return
	}
	messageValue := reflect.ValueOf(message).Elem()
	contextField := messageValue.FieldByName("Context")
	if contextField.CanSet() && contextField.Type().String() == "context.Context" {
		contextField.Set(reflect.ValueOf(request.Context()))
	}
}

func bindJSON(request *http.Request, message interface{}) error {
	binder, ok := message.(BindJSON)
	if !ok {
		return nil
	}
	if !binder.BindJSON() {
		return nil
	}
	if !isPutOrPost(request) {
		return errMethodNotAllowed
	}
	if !hasJSONContent(request) {
		return errUnsupportedMediaType
	}

	return json.NewDecoder(request.Body).Decode(&message)
}
func isPutOrPost(request *http.Request) bool {
	return request.Method == http.MethodPost || request.Method == http.MethodPut
}
func hasJSONContent(request *http.Request) bool {
	return strings.Contains(request.Header.Get("Content-Type"), "/json")
}

func statusCodeFromErrorOrDefault(err error, defaultStatusCode int) (int, error) {
	status, ok := err.(ErrorCode)
	if !ok {
		return defaultStatusCode, err
	}

	code := status.StatusCode()
	if code == 0 {
		return defaultStatusCode, err
	}

	return code, err
}

func sanitize(message interface{}) {
	if sanitizer, isSanitizer := message.(Sanitizer); isSanitizer {
		sanitizer.Sanitize()
	}
}

func validate(message interface{}) error {
	validator, isValidator := message.(Validator)
	if !isValidator {
		return nil
	}

	err := validator.Validate()
	if err == nil {
		return nil
	}

	errs, isErrors := err.(Errors)
	if isErrors && len(errs) == 0 {
		return nil
	}

	diagnosticErrs, isDiagnosticErrors := err.(DiagnosticErrors)
	if isDiagnosticErrors && len(diagnosticErrs) == 0 {
		return nil
	}

	return err
}

func serverError(message interface{}) error {
	if server, isServerError := message.(ServerError); isServerError && server.Error() {
		return internalServerError
	}
	return nil
}

var (
	internalServerError     = errors.New(http.StatusText(http.StatusInternalServerError))
	errUnsupportedMediaType = NewStatusCodeError(http.StatusUnsupportedMediaType)
	errMethodNotAllowed     = NewStatusCodeError(http.StatusMethodNotAllowed)
)

//////////////////////////////////////////////////////////////////////

type StatusCodeError struct {
	statusCode int
}

func NewStatusCodeError(statusCode int) *StatusCodeError {
	return &StatusCodeError{statusCode: statusCode}
}

func (this StatusCodeError) StatusCode() int {
	return this.statusCode
}

func (this StatusCodeError) Error() string {
	return http.StatusText(this.statusCode)
}

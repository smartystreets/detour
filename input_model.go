package detour

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func prepareInputModel(model interface{}, request *http.Request) (statusCode int, err error) {
	if err = bind(request, model); err != nil {
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

func bind(request *http.Request, message interface{}) error {
	if canBindJSON(request, message) {
		if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
			return err
		}
	}

	binder, isBinder := message.(Binder)
	if !isBinder {
		return nil
	}

	err := request.ParseForm()
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

	return err
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

func canBindJSON(request *http.Request, message interface{}) bool {
	if request.Method != http.MethodPost && request.Method != http.MethodPut {
		return false
	}
	if !strings.Contains(request.Header.Get("Content-Type"), "/json") {
		return false
	}
	binder, ok := message.(BindJSON)
	if !ok {
		return false
	}
	if !binder.BindJSON() {
		return false
	}
	return true
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

	return err
}

func serverError(message interface{}) error {
	if server, isServerError := message.(ServerError); isServerError && server.Error() {
		return internalServerError
	}
	return nil
}

var internalServerError = errors.New(http.StatusText(http.StatusInternalServerError))

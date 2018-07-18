package detour

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func prepareInputModel(model interface{}, request *http.Request) (statusCode int, err error) {
	if err = bind(request, model); err != nil {
		return http.StatusBadRequest, err
	}

	sanitize(model)

	if err = validate(model); err != nil {
		return http.StatusUnprocessableEntity, err
	}

	if err = serverError(model); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func bind(request *http.Request, message interface{}) error {
	if isJSON(request) {
		return json.NewDecoder(request.Body).Decode(&message)
	} else if binder, isBinder := message.(Binder); !isBinder {
		return nil
	} else if err := request.ParseForm(); err != nil {
		return err
	} else if err = binder.Bind(request); err == nil {
		return nil
	} else if errs, isErrors := err.(Errors); isErrors && len(errs) == 0 {
		return nil
	} else {
		return err
	}
}

func isJSON(request *http.Request) bool {
	return (request.Method == "POST" || request.Method == "PUT") &&
		strings.Contains(request.Header.Get("Content-Type"), "/json")
}

func sanitize(message interface{}) {
	if sanitizer, isSanitizer := message.(Sanitizer); isSanitizer {
		sanitizer.Sanitize()
	}
}
func validate(message interface{}) error {
	if validator, isValidator := message.(Validator); !isValidator {
		return nil
	} else if err := validator.Validate(); err == nil {
		return nil
	} else if errs, isErrors := err.(Errors); isErrors && len(errs) == 0 {
		return nil
	} else {
		return err
	}
}

func serverError(message interface{}) error {
	if server, isServerError := message.(ServerError); isServerError && server.Error() {
		return internalServerError
	}
	return nil
}

var internalServerError = errors.New(http.StatusText(http.StatusInternalServerError))

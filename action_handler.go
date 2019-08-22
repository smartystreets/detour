package detour

import (
	"net/http"
)

type actionHandler struct {
	controller            monadicAction
	generateNewInputModel createModel
}

// Install merely allows *actionHandler to implement a non-public/internal, company-specific interface.
func (this *actionHandler) Install(http.Handler) {}

func (this *actionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	model := this.generateNewInputModel()
	status, err := prepareInputModel(model, request)
	result := this.determineResult(model, status, err)
	result.Render(response, request)
}

func (this *actionHandler) determineResult(model interface{}, status int, err error) Renderer {
	if err != nil {
		return inputModelErrorResult(status, err)
	} else {
		return this.controllerActionResult(model)
	}
}

func inputModelErrorResult(code int, err error) Renderer {
	_, isErrors := err.(Errors)
	if isErrors {
		return &JSONResult{StatusCode: code, Content: err}
	}

	_, isDiagnosticErr := err.(*DiagnosticError)
	if isDiagnosticErr {
		return &DiagnosticResult{StatusCode: code, Message: err.Error()}
	}

	_, isDiagnosticErrors := err.(DiagnosticErrors)
	if isDiagnosticErrors {
		return &DiagnosticResult{StatusCode: code, Message: http.StatusText(code) + "\n\n" + err.Error()}
	}

	return &StatusCodeResult{StatusCode: code, Message: err.Error()}
}

func (this *actionHandler) controllerActionResult(model interface{}) Renderer {
	if result := this.controller(model); result != nil {
		return result
	} else {
		return nopResult{}
	}
}

type nopResult struct{}

func (nopResult) Render(http.ResponseWriter, *http.Request) {}

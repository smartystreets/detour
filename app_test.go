package detour

import (
	"net/http"
	"strings"
)

///////////////////////////////////////////////////////////////

type Controller struct{}

func (this *Controller) HandleBasicInputModel(model *BlankBasicInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleBindingInputModel(model *BindingInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleBindingFailsInputModel(model *BindingFailsInputModel) Renderer {
	panic("We shouldn't reach this point because the binding failed.")
}
func (this *Controller) HandleBindingFailsCustomStatusCodeInputModel(model *BindingFailsWithCustomStatusCodeInputModel) Renderer {
	panic("We shouldn't reach this point because the binding failed.")
}
func (this *Controller) HandleBindingEmptyErrorsInputModel(model *BindingEmptyErrorsInputModel) Renderer {
	return &ControllerResponse{}
}
func (this *Controller) HandleBindingFailsInputModelWithDiagnostics(model *BindingFailsInputModelWithDiagnostics) Renderer {
	panic("We shouldn't reach this point because the binding failed.")
}
func (this *Controller) HandleSanitizingInputModel(model *SanitizingInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleValidatingInputModel(model *ValidatingInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleValidatingEmptyErrors(model *ValidatingEmptyErrorsInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleValidatingFailsInputModel(model *ValidatingFailsInputModel) Renderer {
	panic("We shouldn't reach this point because the validation failed.")
}
func (this *Controller) HandleValidatingFailsWithCustomStatusCode(model *ValidatingFailsWithCustomStatusCodeInputModel) Renderer {
	panic("We shouldn't reach this point because the validation failed.")
}
func (this *Controller) HandleFinalError(model *FinalErrorInputModel) Renderer {
	panic("We should't reach this point because the Error method returned true.")
}
func (this *Controller) HandleNoFinalError(model *NoFinalErrorInputModel) Renderer {
	return nil
}
func (this *Controller) HandleNilResponseInputModel(model *NilResponseInputModel) Renderer {
	return nil
}
func (this *Controller) HandleEmptyInputModel() Renderer {
	return nil
}

///////////////////////////////////////////////////////////////

type BlankBasicInputModel struct {
	Content string
}

func NewBlankBasicInputModel() interface{} {
	return &BlankBasicInputModel{Content: "BasicInputModel"}
}

/////

type BindingInputModel struct {
	Content string
}

func (this *BindingInputModel) Bind(request *http.Request) error {
	this.Content = request.Form.Get("binding")
	return nil
}

/////

type BindingEmptyErrorsInputModel struct {
}

func (this *BindingEmptyErrorsInputModel) Bind(request *http.Request) error {
	var errors Errors
	return errors
}

/////

type BindingFailsInputModel struct{}

func (this *BindingFailsInputModel) Bind(request *http.Request) error {
	var errors Errors
	errors = errors.Append(NewBindingValidationError("BindingFailsInputModel"))
	return errors
}

/////

type BindingFailsWithCustomStatusCodeInputModel struct{}

func (this *BindingFailsWithCustomStatusCodeInputModel) Bind(request *http.Request) error {
	var errors Errors
	errors = errors.Append(&InputError{HTTPStatusCode: http.StatusTeapot})
	return errors
}

/////

type BindingFailsInputModelWithDiagnostics struct{}

func (this *BindingFailsInputModelWithDiagnostics) Bind(request *http.Request) error {
	return NewDiagnosticError("BindingFailsInputModel")
}

/////

type SanitizingInputModel struct {
	Content string
}

func (this *SanitizingInputModel) Bind(request *http.Request) error {
	this.Content = request.Form.Get("binding")
	return nil
}

func (this *SanitizingInputModel) Sanitize() {
	this.Content = strings.ToUpper(strings.Replace(this.Content, "Binding", "Sanitizing", 1))
}

/////

type ValidatingInputModel struct {
	Content string
}

func (this *ValidatingInputModel) Bind(request *http.Request) error {
	return nil
}

func (this *ValidatingInputModel) Validate() error {
	this.Content = "ValidatingInputModel"
	return nil
}

/////

type ValidatingFailsInputModel struct{}

func (this *ValidatingFailsInputModel) Validate() error {
	var errors Errors
	errors = errors.Append(NewBindingValidationError("ValidatingFailsInputModel"))
	return errors
}

/////

type ValidatingFailsWithCustomStatusCodeInputModel struct{}

func (this *ValidatingFailsWithCustomStatusCodeInputModel) Validate() error {
	return &DiagnosticError{HTTPStatusCode: http.StatusTeapot}
}

/////

type ValidatingEmptyErrorsInputModel struct{ Content string }

func (this *ValidatingEmptyErrorsInputModel) Validate() error {
	this.Content = "ValidatingEmptyErrorsInputModel"
	var errors Errors
	return errors
}

/////

type FinalErrorInputModel struct{}

func (this *FinalErrorInputModel) Error() bool { return true }

/////

type NoFinalErrorInputModel struct{}

func (this *NoFinalErrorInputModel) Error() bool { return false }

/////

type NilResponseInputModel struct{}

////////////////////////////////////////////////////////////////

type ControllerResponse struct {
	Body string
}

func (this *ControllerResponse) Render(response http.ResponseWriter, request *http.Request) {
	http.Error(response, "Just handled: "+this.Body, http.StatusOK)
}

////////////////////////////////////////////////////////////////

type BindingValidationError struct {
	Problem string
}

func NewBindingValidationError(problem string) *BindingValidationError {
	return &BindingValidationError{Problem: problem}
}

func (this *BindingValidationError) Error() string {
	return this.Problem
}

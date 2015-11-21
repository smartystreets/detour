package binding

import (
	"fmt"
	"net/http"
)

///////////////////////////////////////////////////////////////

type Application struct{}

func (this *Application) HandleBasicInputModel(model interface{}) http.Handler {
	return &ApplicationResponse{Body: model.(*BlankBasicInputModel).Content}
}
func (this *Application) HandleBindingInputModel(model interface{}) http.Handler {
	return &ApplicationResponse{Body: model.(*BindingInputModel).Content}
}
func (this *Application) HandleBindingFailsInputModel(model interface{}) http.Handler {
	panic("We shouldn't reach this point because the binding failed.")
}
func (this *Application) HandleValidatingInputModel(model interface{}) http.Handler {
	return &ApplicationResponse{Body: model.(*ValidatingInputModel).Content}
}
func (this *Application) HandleValidatingEmptyErrors(model interface{}) http.Handler {
	return &ApplicationResponse{Body: model.(*ValidatingEmptyErrorsInputModel).Content}
}
func (this *Application) HandleValidatingFailsInputModel(model interface{}) http.Handler {
	panic("We shouldn't reach this point because the validation failed.")
}
func (this *Application) HandleTranslatingInputModel(model interface{}) http.Handler {
	return &ApplicationResponse{Body: model.(string)}
}
func (this *Application) HandleNilResponseInputModel(model interface{}) http.Handler {
	return nil
}

/////

type GenericHandler struct{}

func (this *GenericHandler) Handle(response http.ResponseWriter, request *http.Request, model interface{}) {
	fmt.Fprint(response, "Just handled: ", model.(*GenericHandlerInputModel).Content)
}

///////////////////////////////////////////////////////////////

type BlankBasicInputModel struct {
	Content string
}

func NewBlankBasicInputModel() interface{} {
	return &BlankBasicInputModel{
		Content: "BasicInputModel",
	}
}

/////

type BindingInputModel struct {
	Content string
}

func NewBindingInputModel() interface{} {
	return &BindingInputModel{}
}

func (this *BindingInputModel) Bind(request *http.Request) error {
	this.Content = request.Form.Get("binding")
	return nil
}

/////

type BindingFailsInputModel struct{}

func NewBindingFailsInputModel() interface{} {
	return &BindingFailsInputModel{}
}

func (this *BindingFailsInputModel) Bind(request *http.Request) error {
	return NewBindingValidationError("BindingFailsInputModel")
}

/////

type ValidatingInputModel struct {
	Content string
}

func NewValidatingInputModel() interface{} {
	return &ValidatingInputModel{}
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

func NewValidatingFailsInputModel() interface{} {
	return &ValidatingFailsInputModel{}
}

func (this *ValidatingFailsInputModel) Validate() error {
	return NewBindingValidationError("ValidatingFailsInputModel")
}

/////

type ValidatingEmptyErrorsInputModel struct{ Content string }

func NewValidatingEmptyInputModel() interface{} {
	return &ValidatingEmptyErrorsInputModel{}
}
func (this *ValidatingEmptyErrorsInputModel) Validate() error {
	this.Content = "ValidatingEmptyErrorsInputModel"
	var errors ValidationErrors
	return errors
}

/////

type TranslatingInputModel struct{}

func NewTranslatingInputModel() interface{} {
	return &TranslatingInputModel{}
}

func (this *TranslatingInputModel) Translate() interface{} {
	return "TranslatingInputModel"
}

/////

type NilResponseInputModel struct{}

func NewNilResponseInputModel() interface{} {
	return &NilResponseInputModel{}
}

/////

type GenericHandlerInputModel struct {
	Content string
}

func NewGenericHandlerInputModel() interface{} {
	return &GenericHandlerInputModel{
		Content: "GenericHandlerInputModel",
	}
}

////////////////////////////////////////////////////////////////

type ApplicationResponse struct {
	Body string
}

func (this *ApplicationResponse) ServeHTTP(response http.ResponseWriter, request *http.Request) {
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

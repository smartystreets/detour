package binding

import "net/http"

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
func (this *Controller) HandleValidatingInputModel(model *ValidatingInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleValidatingEmptyErrors(model *ValidatingEmptyErrorsInputModel) Renderer {
	return &ControllerResponse{Body: model.Content}
}
func (this *Controller) HandleValidatingFailsInputModel(model *ValidatingFailsInputModel) Renderer {
	panic("We shouldn't reach this point because the validation failed.")
}
func (this *Controller) HandleNilResponseInputModel(model *NilResponseInputModel) Renderer {
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

type BindingFailsInputModel struct{}

func (this *BindingFailsInputModel) Bind(request *http.Request) error {
	return NewBindingValidationError("BindingFailsInputModel")
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
	return NewBindingValidationError("ValidatingFailsInputModel")
}

/////

type ValidatingEmptyErrorsInputModel struct{ Content string }

func (this *ValidatingEmptyErrorsInputModel) Validate() error {
	this.Content = "ValidatingEmptyErrorsInputModel"
	var errors ValidationErrors
	return errors
}

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

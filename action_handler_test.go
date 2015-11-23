package binding

import (
	"net/http"
	"net/http/httptest"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ModelBinderFixture struct {
	*gunit.Fixture

	controller *Controller
	request    *http.Request
	response   *httptest.ResponseRecorder
}

func (this *ModelBinderFixture) Setup() {
	this.controller = &Controller{}
	this.request, _ = http.NewRequest("GET", "/?binding=BindingInputModel", nil)
	this.response = httptest.NewRecorder()
}

func (this *ModelBinderFixture) TestBasicInputModelProvidedToApplication__HTTP200() {
	binder := withFactory(this.controller.HandleBasicInputModel, NewBlankBasicInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BasicInputModel")
}

func (this *ModelBinderFixture) TestBindsModelForApplication__HTTP200() {
	binder := New(this.controller.HandleBindingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesError__HTTP400() {
	binder := New(this.controller.HandleBindingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "BindingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelForApplication__HTTP200() {
	binder := New(this.controller.HandleValidatingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelAndHandlesError__HTTP422() {
	binder := New(this.controller.HandleValidatingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 422)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "ValidatingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelEmptyValidationErrors__HTTP200() {
	binder := New(this.controller.HandleValidatingEmptyErrors)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingEmptyErrorsInputModel")
}

func (this *ModelBinderFixture) TestNilResponseFromApplication__HTTP200() {
	binder := New(this.controller.HandleNilResponseInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.BeBlank)
}

////////////////////////////////////////////////////////////

func (this *ModelBinderFixture) TestModelParsingFromCallback() {
	this.assertPanic(0)                                   // not a method
	this.assertPanic(func() {})                           // no input
	this.assertPanic(func(int) {})                        // not a pointer
	this.So(func() { parseModelType(func(*BlankBasicInputModel) {} ) }, should.Panic) // doesn't return a Renderer
	this.So(func() { parseModelType(func(*BlankBasicInputModel) Renderer { return nil }) }, should.NotPanic)
}
func (this *ModelBinderFixture) assertPanic(callback interface{}) {
	this.So(func() { parseModelType(callback) }, should.Panic)
}

func (this *ModelBinderFixture) TestModelBinding() {
	binder := New(func(input *BindingInputModel) Renderer {
		return &ControllerResponse{Body: input.Content}
	})

	binder.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

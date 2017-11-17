package detour

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestModelBinderFixture(t *testing.T) {
	gunit.Run(new(ModelBinderFixture), t)
}

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

func (this *ModelBinderFixture) TestNoInputModelProvidedToApplication__HTTP200() {
	binder := New(this.controller.HandleEmptyInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
}

func (this *ModelBinderFixture) TestBindsModelForApplication__HTTP200() {
	binder := New(this.controller.HandleBindingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

func (this *ModelBinderFixture) TestBindsFormParseFails__HTTP400() {
	binder := New(this.controller.HandleBindingInputModel)
	this.request.Method = "PUT"
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusBadRequest)
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesError__HTTP400_JSONResponse() {
	binder := New(this.controller.HandleBindingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.HeaderMap.Get(contentTypeHeader), should.Equal, jsonContentType)
	this.So(this.response.Body.String(), should.EqualTrimSpace, `[{"Problem":"BindingFailsInputModel"}]`)
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesNilErrors() {
	binder := New(this.controller.HandleBindingEmptyErrorsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
}

func (this *ModelBinderFixture) TestBindModelAndHandleError__HTTP400_DiagnosticsResponse() {
	binder := New(this.controller.HandleBindingFailsInputModelWithDiagnostics)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.HeaderMap.Get(contentTypeHeader), should.Equal, plaintextContentType)
	this.So(this.response.Body.String(), should.ContainSubstring, `400 BindingFailsInputModel`)
	this.So(this.response.Body.String(), should.ContainSubstring, "Raw Request:")
	this.So(this.response.Body.String(), should.ContainSubstring, "/?binding=BindingInputModel") // from the URL
	this.So(this.response.Body.String(), should.ContainSubstring, `---- DISCLAIMER ----`)
}

func (this *ModelBinderFixture) TestSanitizesModelIfAvailable() {
	sanitizer := New(this.controller.HandleSanitizingInputModel)
	sanitizer.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: SANITIZINGINPUTMODEL")
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
	this.So(this.response.HeaderMap.Get(contentTypeHeader), should.Equal, jsonContentType)
	this.So(this.response.Body.String(), should.EqualTrimSpace, `[{"Problem":"ValidatingFailsInputModel"}]`)
}

func (this *ModelBinderFixture) TestValidatesModelEmptyValidationErrors__HTTP200() {
	binder := New(this.controller.HandleValidatingEmptyErrors)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingEmptyErrorsInputModel")
}

func (this *ModelBinderFixture) TestFinalErrorCondition__HTTP500() {
	action := New(this.controller.HandleFinalError)
	action.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 500)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Internal Server Error")
}

func (this *ModelBinderFixture) TestNoFinalErrorCondition__HTTP200() {
	action := New(this.controller.HandleNoFinalError)
	action.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
}

func (this *ModelBinderFixture) TestNilResponseFromApplication__HTTP200() {
	binder := New(this.controller.HandleNilResponseInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.BeBlank)
}

////////////////////////////////////////////////////////////

func (this *ModelBinderFixture) TestModelParsingFromCallback() {
	this.assertPanic(0)                                                                              // not a method
	this.assertPanic(func() {})                                                                      // no input
	this.assertPanic(func(int) {})                                                                   // not a pointer
	this.assertPanic(func(*int, *int) {})                                                            // not a pointer
	this.So(func() { identifyInputModelArgumentType(func(*BlankBasicInputModel) {}) }, should.Panic) // doesn't return a Renderer
	this.So(func() { identifyInputModelArgumentType(func(*BlankBasicInputModel) Renderer { return nil }) }, should.NotPanic)
}
func (this *ModelBinderFixture) assertPanic(callback interface{}) {
	this.So(func() { identifyInputModelArgumentType(callback) }, should.Panic)
}

func (this *ModelBinderFixture) TestModelBinding() {
	binder := New(func(input *BindingInputModel) Renderer {
		return &ControllerResponse{Body: input.Content}
	})

	binder.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

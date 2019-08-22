package detour

import (
	"net/http"
	"net/http/httptest"
	"strings"
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
	this.request = httptest.NewRequest("GET", "/?binding=BindingInputModel", nil)
	this.response = httptest.NewRecorder()
}

func (this *ModelBinderFixture) TestFromFactory_IncorrectInputModelType__Panic() {
	wrongInputModelType := func() interface{} { return "wrong type" }
	action := func() { NewFromFactory(wrongInputModelType, this.controller.HandleBasicInputModel) }
	this.So(action, should.PanicWith,
		"Controller requires input model of type: [*detour.BlankBasicInputModel] " +
		"Factory function provided input model of type: [string]")
}
func (this *ModelBinderFixture) TestFromFactory_ControllerWithNoInputModel__Panic() {
	action := func() { NewFromFactory(NewBlankBasicInputModel, this.controller.HandleNoInputModel) }
	this.So(action, should.Panic)
}
func (this *ModelBinderFixture) TestFromFactory_BasicInputModelProvidedToApplication__HTTP200() {
	binder := NewFromFactory(NewBlankBasicInputModel, this.controller.HandleBasicInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BasicInputModel")
}

func (this *ModelBinderFixture) TestNoInputModelProvidedToApplication__HTTP200() {
	binder := New(this.controller.HandleNoInputModel)
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
	this.request = httptest.NewRequest("GET", "/?asdf=%%%%%", nil)
	binder := New(this.controller.HandleBindingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusBadRequest)
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesError__HTTP400_JSONResponse() {
	binder := New(this.controller.HandleBindingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Result().Header.Get(contentTypeHeader), should.Equal, jsonContentType)
	this.So(this.response.Body.String(), should.EqualTrimSpace, `[{"Problem":"BindingFailsInputModel"}]`)
}

func (this *ModelBinderFixture) TestBindModelError__CustomStatusCode() {
	binder := New(this.controller.HandleBindingFailsCustomStatusCodeInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusTeapot)
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
	this.So(this.response.Result().Header.Get(contentTypeHeader), should.Equal, plaintextContentType)
	this.So(this.response.Body.String(), should.ContainSubstring, "400 BindingFailsInputModel")
	this.So(this.response.Body.String(), should.ContainSubstring, "Raw Request:")
	this.So(this.response.Body.String(), should.ContainSubstring, "/?binding=BindingInputModel") // from the URL
	this.So(this.response.Body.String(), should.ContainSubstring, "---- DISCLAIMER ----")
}

func (this *ModelBinderFixture) TestBindModelAndHandleError__HTTP400_DiagnosticErrorsResponse() {
	binder := New(this.controller.HandleBindingFailsInputModelWithDiagnosticErrors)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Result().Header.Get(contentTypeHeader), should.Equal, plaintextContentType)
	this.So(this.response.Body.String(), should.ContainSubstring, "400 Bad Request")
	this.So(this.response.Body.String(), should.ContainSubstring, "BindingFailsInputModel")
	this.So(this.response.Body.String(), should.ContainSubstring, "Raw Request:")
	this.So(this.response.Body.String(), should.ContainSubstring, "/?binding=BindingInputModel") // from the URL
	this.So(this.response.Body.String(), should.ContainSubstring, "---- DISCLAIMER ----")
}

func (this *ModelBinderFixture) TestBindFromJSONRequiresInputModelToImplementJSONMarkerInterface() {
	this.request = httptest.NewRequest("POST", "/", strings.NewReader(`{"content": "Hello, World!"}`))
	this.request.Header.Set("Content-Type", "application/json")
	binder := New(this.controller.HandleFailedBindingFromJSON)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.NotContainSubstring, "Hello, World!")
}

func (this *ModelBinderFixture) TestBindFromJSONPost() {
	this.request = httptest.NewRequest("POST", "/", strings.NewReader(`{"content": "Hello, World!"}`))
	this.request.Header.Set("Content-Type", "application/json")
	this.request.Header.Set("binding", " (from the header)")
	binder := New(this.controller.HandleBindingFromJSON)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.ContainSubstring, "Hello, World! (from the header)")
}
func (this *ModelBinderFixture) TestBindFromJSONPut() {
	this.request = httptest.NewRequest("PUT", "/", strings.NewReader(`{"content": "Hello, World!"}`))
	this.request.Header.Set("Content-Type", "application/json")
	binder := New(this.controller.HandleBindingFromJSON)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.ContainSubstring, "Hello, World!")
}
func (this *ModelBinderFixture) TestBindFromJSON_Malformed() {
	this.request = httptest.NewRequest("PUT", "/", strings.NewReader(`{I can haz JSONs}`))
	this.request.Header.Set("Content-Type", "application/json")
	binder := New(this.controller.HandleBindingFromJSON)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Body.String(), should.ContainSubstring, "invalid character")
}

func (this *ModelBinderFixture) TestSanitizesModelIfAvailable() {
	this.request.Header.Set("binding", "hello")
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
	this.So(this.response.Result().Header.Get(contentTypeHeader), should.Equal, jsonContentType)
	this.So(this.response.Body.String(), should.EqualTrimSpace, `[{"Problem":"ValidatingFailsInputModel"}]`)
}

func (this *ModelBinderFixture) TestValidatesModelEmptyValidationErrors__HTTP200() {
	binder := New(this.controller.HandleValidatingEmptyErrors)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingEmptyErrorsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelCustomStatusCodeError__HTTP418() {
	binder := New(this.controller.HandleValidatingFailsWithCustomStatusCode)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusTeapot)
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

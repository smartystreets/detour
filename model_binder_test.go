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
	binder := TypedFactory(this.controller.HandleBasicInputModel, NewBlankBasicInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BasicInputModel")
}

func (this *ModelBinderFixture) TestBindsModelForApplication__HTTP200() {
	binder := Typed(this.controller.HandleBindingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesError__HTTP400() {
	binder := Typed(this.controller.HandleBindingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "BindingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelForApplication__HTTP200() {
	binder := Typed(this.controller.HandleValidatingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelAndHandlesError__HTTP422() {
	binder := Typed(this.controller.HandleValidatingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 422)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "ValidatingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelEmptyValidationErrors__HTTP200() {
	binder := Typed(this.controller.HandleValidatingEmptyErrors)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingEmptyErrorsInputModel")
}

func (this *ModelBinderFixture) TestNilResponseFromApplication__HTTP200() {
	binder := Typed(this.controller.HandleNilResponseInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.BeBlank)
}

////////////////////////////////////////////////////////////

//func (this *ModelBinderFixture) TestInputModelParsingFromCallback() {
//	this.assertPanic(0)                                                                 // not a method
//	this.assertPanic(func(int) {})                                                      // too few arguments (need 3 arguments)
//	this.assertPanic(func(int, int, int) {})                                            // bad first argument (not http.ResponseWriter)
//	this.assertPanic(func(http.ResponseWriter, int, int) {})                            // bad second argument (not *http.Request)
//	this.assertPanic(func(http.ResponseWriter, *http.Request, BlankBasicInputModel) {}) // bad third argument (not a pointer)
//	this.So(func() { parseInputModelType(func(http.ResponseWriter, *http.Request, *BlankBasicInputModel) {}) }, should.NotPanic)
//}
//func (this *ModelBinderFixture) assertPanic(callback interface{}) {
//	this.So(func() { parseInputModelType(callback) }, should.Panic)
//}

//func (this *ModelBinderFixture) TestModelBinding() {
//	binder := Typed(func(input *BindingInputModel) Renderer { fmt.Fprintf(w, input.Content) })
//
//	binder.ServeHTTP(this.response, this.request)
//
//	this.So(this.response.Code, should.Equal, 200)
//	this.So(this.response.Body.String(), should.Equal, "BindingInputModel")
//}

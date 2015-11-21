package binding

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ModelBinderFixture struct {
	*gunit.Fixture

	app      *Application
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *ModelBinderFixture) Setup() {
	this.app = &Application{}
	this.request, _ = http.NewRequest("GET", "/?binding=BindingInputModel", nil)
	this.response = httptest.NewRecorder()
}

func (this *ModelBinderFixture) TestBasicInputModelProvidedToApplication__HTTP200() {
	binder := DomainBinder(NewBlankBasicInputModel, this.app.HandleBasicInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BasicInputModel")
}

func (this *ModelBinderFixture) TestBindsModelForApplication__HTTP200() {
	binder := DomainBinder(NewBindingInputModel, this.app.HandleBindingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesError__HTTP400() {
	binder := DomainBinder(NewBindingFailsInputModel, this.app.HandleBindingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "BindingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelForApplication__HTTP200() {
	binder := DomainBinder(NewValidatingInputModel, this.app.HandleValidatingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelAndHandlesError__HTTP422() {
	binder := DomainBinder(NewValidatingFailsInputModel, this.app.HandleValidatingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 422)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "ValidatingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelEmptyValidationErrors__HTTP200() {
	binder := DomainBinder(NewValidatingEmptyInputModel, this.app.HandleValidatingEmptyErrors)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingEmptyErrorsInputModel")
}

func (this *ModelBinderFixture) TestTranslatesModelForApplication__HTTP200() {
	binder := DomainBinder(NewTranslatingInputModel, this.app.HandleTranslatingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: TranslatingInputModel")
}

func (this *ModelBinderFixture) TestNilResponseFromApplication__HTTP200() {
	binder := DomainBinder(NewNilResponseInputModel, this.app.HandleNilResponseInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.BeBlank)
}

func (this *ModelBinderFixture) TestGenericHandlerAsApplication__HTTP200() {
	binder := ControllerBinder(NewGenericHandlerInputModel, (&GenericHandler{}).Handle)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: GenericHandlerInputModel")
}

////////////////////////////////////////////////////////////

func (this *ModelBinderFixture) TestInputModelParsingFromCallback() {
	this.assertPanic(0)                                                                 // not a method
	this.assertPanic(func(int) {})                                                      // too few arguments (need 3 arguments)
	this.assertPanic(func(int, int, int) {})                                            // bad first argument (not http.ResponseWriter)
	this.assertPanic(func(http.ResponseWriter, int, int) {})                            // bad second argument (not *http.Request)
	this.assertPanic(func(http.ResponseWriter, *http.Request, BlankBasicInputModel) {}) // bad third argument (not a pointer)
	this.So(func() { parseInputModelType(func(http.ResponseWriter, *http.Request, *BlankBasicInputModel) {}) }, should.NotPanic)
}
func (this *ModelBinderFixture) assertPanic(callback interface{}) {
	this.So(func() { parseInputModelType(callback) }, should.Panic)
}

func (this *ModelBinderFixture) TestModelBinding() {
	binder := DefaultBinder(func(w http.ResponseWriter, r *http.Request, input *BindingInputModel) {
		fmt.Fprintf(w, input.Content)
	})

	binder.ServeHTTP(this.response, this.request)

	this.So(this.response.Code, should.Equal, 200)
	this.So(this.response.Body.String(), should.Equal, "BindingInputModel")
}

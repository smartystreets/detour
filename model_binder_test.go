package binding

import (
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
	binder := NewDomainModelBinder(NewBlankBasicInputModel, this.app.HandleBasicInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BasicInputModel")
}

func (this *ModelBinderFixture) TestBindsModelForApplication__HTTP200() {
	binder := NewDomainModelBinder(NewBindingInputModel, this.app.HandleBindingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: BindingInputModel")
}

func (this *ModelBinderFixture) TestBindsModelAndHandlesError__HTTP400() {
	binder := NewDomainModelBinder(NewBindingFailsInputModel, this.app.HandleBindingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 400)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "BindingFailsInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelForApplication__HTTP200() {
	binder := NewDomainModelBinder(NewValidatingInputModel, this.app.HandleValidatingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: ValidatingInputModel")
}

func (this *ModelBinderFixture) TestValidatesModelAndHandlesError__HTTP422() {
	binder := NewDomainModelBinder(NewValidatingFailsInputModel, this.app.HandleValidatingFailsInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, 422)
	this.So(this.response.Header().Get("Content-Type"), should.Equal, "application/json")
	this.So(this.response.Body.String(), should.ContainSubstring, "ValidatingFailsInputModel")
}

func (this *ModelBinderFixture) TestTranslatesModelForApplication__HTTP200() {
	binder := NewDomainModelBinder(NewTranslatingInputModel, this.app.HandleTranslatingInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: TranslatingInputModel")
}

func (this *ModelBinderFixture) TestNilResponseFromApplication__HTTP200() {
	binder := NewDomainModelBinder(NewNilResponseInputModel, this.app.HandleNilResponseInputModel)
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.BeBlank)
}

func (this *ModelBinderFixture) TestGenericHandlerAsApplication__HTTP200() {
	binder := NewModelBinderHandler(NewGenericHandlerInputModel, &GenericHandler{})
	binder.ServeHTTP(this.response, this.request)
	this.So(this.response.Code, should.Equal, http.StatusOK)
	this.So(this.response.Body.String(), should.EqualTrimSpace, "Just handled: GenericHandlerInputModel")
}

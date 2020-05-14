package detour

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDetourFixture(t *testing.T) {
	gunit.Run(new(DetourFixture), t)
}

type DetourFixture struct {
	*gunit.Fixture

	request  *http.Request
	response *httptest.ResponseRecorder
	model    *TestModel

	handledContext  context.Context
	handledMessages []interface{}
}

func (this *DetourFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.response = httptest.NewRecorder()
	this.model = &TestModel{}
}

func (this *DetourFixture) Handle(ctx context.Context, messages ...interface{}) {
	this.handledContext = ctx
	this.handledMessages = messages
}
func (this *DetourFixture) Model() Detour {
	return this.model
}

func (this *DetourFixture) Test_ReturningRendererFromBindShortCircuitsHandler() {
	New(this.Model, this).ServeHTTP(this.response, this.request)

	this.So(this.model.boundRequest, should.Equal, this.request)
	this.So(this.handledContext, should.Equal, this.request.Context())
	this.So(this.handledMessages, should.Resemble, []interface{}{1, 2, 3})
	this.So(this.response.Result().StatusCode, should.Equal, http.StatusTeapot)
	this.So(this.model.renderedRequest, should.Equal, this.request)
}

///////////////////////////////////////////////////////////////

type TestModel struct {
	boundRequest    *http.Request
	renderedRequest *http.Request
}

func (this *TestModel) Bind(request *http.Request) (messages []interface{}) {
	this.boundRequest = request
	return []interface{}{1, 2, 3}
}

func (this *TestModel) Render(response http.ResponseWriter, request *http.Request) {
	this.renderedRequest = request
	response.WriteHeader(http.StatusTeapot)
}

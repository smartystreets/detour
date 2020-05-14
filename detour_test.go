package detour

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/detour/v3/render"
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
func (this *DetourFixture) AssertHandleNotInvoked() {
	this.So(this.handledContext, should.BeNil)
	this.So(this.handledMessages, should.BeNil)
}
func (this *DetourFixture) AssertHandled(expectedContext context.Context, expectedMessages ...interface{}) {
	this.So(this.handledContext, should.Resemble, expectedContext)
	this.So(this.handledMessages, should.Resemble, expectedMessages)
}
func (this *DetourFixture) AssertStatusCode(expectedStatusCode int) bool {
	return this.So(this.response.Result().StatusCode, should.Equal, expectedStatusCode)
}
func (this *DetourFixture) Model() Detour { return this.model }

func (this *DetourFixture) Test_ReturningRendererFromBindShortCircuitsHandler() {
	this.model.bindRenderer = render.StatusCodeResult{StatusCode: http.StatusTeapot}

	New(this.Model, this).ServeHTTP(this.response, this.request)

	this.So(this.model.boundRequest, should.Equal, this.request)
	this.AssertHandleNotInvoked()
	this.AssertStatusCode(http.StatusTeapot)
}
func (this *DetourFixture) Test_ReturningNilFromBindAllowsHandlerToBeInvoked() {
	this.model.bindRenderer = nil
	this.model.handleRenderer = render.StatusCodeResult{StatusCode: http.StatusCreated}
	this.model.messages = []interface{}{1, 2, 3}

	New(this.Model, this).ServeHTTP(this.response, this.request)

	this.So(this.model.boundRequest, should.Equal, this.request)
	this.AssertHandled(this.request.Context(), 1, 2, 3)
	this.AssertStatusCode(http.StatusCreated)
}
func (this *DetourFixture) Test_ReturningNilFromBindAndHandler_HTTP200() {
	this.model.bindRenderer = nil
	this.model.handleRenderer = nil
	this.model.messages = []interface{}{1, 2, 3}

	New(this.Model, this).ServeHTTP(this.response, this.request)

	this.So(this.model.boundRequest, should.Equal, this.request)
	this.AssertHandled(this.request.Context(), 1, 2, 3)
	this.AssertStatusCode(http.StatusOK)
}

///////////////////////////////////////////////////////////////

type TestModel struct {
	boundRequest   *http.Request
	bindRenderer   render.Renderer
	messages       []interface{}
	handleRenderer render.Renderer
}

func (this *TestModel) Bind(request *http.Request) render.Renderer {
	this.boundRequest = request
	return this.bindRenderer
}

func (this *TestModel) MessagesToHandle() []interface{} {
	return this.messages
}

func (this *TestModel) Render() render.Renderer {
	return this.handleRenderer
}

type StringRenderer string

func (this StringRenderer) Render(response http.ResponseWriter, request *http.Request) {
	_, _ = io.WriteString(response, string(this))
}

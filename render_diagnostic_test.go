package detour

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDiagnosticResultFixture(t *testing.T) {
	gunit.Run(new(DiagnosticResultFixture), t)
}

type DiagnosticResultFixture struct {
	*gunit.Fixture
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *DiagnosticResultFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/hello-world", strings.NewReader("Hello, World!"))
	this.response = httptest.NewRecorder()
}

func (this *DiagnosticResultFixture) TestRenderedResponse_DefaultValues() {
	result := &DiagnosticResult{
		StatusCode: 200,
		Message:    "OK",
	}
	result.Render(this.response, this.request)
	body := this.response.Body.String()
	this.So(body, should.StartWith, "200 OK")
	this.So(body, should.ContainSubstring, "GET /hello-world HTTP/1.1")
	this.So(body, should.ContainSubstring, "Host: example.com")
	this.So(body, should.NotContainSubstring, "Hello, World!")
}

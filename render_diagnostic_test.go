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
	this.So(this.response.Body.String(), should.StartWith, "200 OK")
	this.So(this.response.Body.String(), should.ContainSubstring, "GET /hello-world HTTP/1.1")
	this.So(this.response.Body.String(), should.ContainSubstring, "Host: example.com")
	this.So(this.response.Body.String(), should.NotContainSubstring, "Hello, World!")
}

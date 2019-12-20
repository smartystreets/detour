package detour

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/smartystreets/assertions/should"
)

func (this *ResultFixture) TestRenderedResponse() {
	this.request = httptest.NewRequest("GET", "/hello-world", strings.NewReader("Hello, World!"))
	result := &DiagnosticResult{
		StatusCode: 200,
		Message:    "OK",
		Header:     http.Header{"X-Hello": []string{"World"}},
	}
	result.Render(this.response, this.request)
	rawBody, err := ioutil.ReadAll(this.response.Result().Body)
	body := string(rawBody)
	this.So(err, should.BeNil)
	this.So(body, should.StartWith, "200 OK")
	this.So(body, should.ContainSubstring, "> GET /hello-world HTTP/1.1")
	this.So(body, should.ContainSubstring, "> Host: example.com")
	this.So(strings.Count(body, ">"), should.Equal, 4)
	this.So(body, should.NotContainSubstring, "Hello, World!")
	this.assertHasHeader("X-Hello", "World")
}

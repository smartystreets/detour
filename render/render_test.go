package render

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestResultFixture(t *testing.T) {
	gunit.Run(new(ResultFixture), t)
}

type ResultFixture struct {
	*gunit.Fixture

	response *httptest.ResponseRecorder
	request  *http.Request
}

func (this *ResultFixture) Setup() {
	this.response = httptest.NewRecorder()
	this.request = httptest.NewRequest("GET", "/", nil)
}

func (this *ResultFixture) setRequestURLCallback(value string) {
	query := this.request.URL.Query()
	query.Set("callback", value)
	this.request.URL.RawQuery = query.Encode()
}
func (this *ResultFixture) render(result Renderer) {
	result.Render(this.response, this.request)
}
func (this *ResultFixture) assertStatusCode(expected int) {
	this.So(this.response.Result().StatusCode, should.Equal, expected)
}
func (this *ResultFixture) assertContent(expected string) {
	body, _ := ioutil.ReadAll(this.response.Result().Body)
	this.So(string(body), should.EqualTrimSpace, expected)
}
func (this *ResultFixture) assertNoHeaders() {
	this.So(this.response.Result().Header, should.HaveLength, 0)
}
func (this *ResultFixture) assertHasHeader(key, value string) {
	this.So(this.response.Result().Header, should.ContainKey, key)
	this.So(this.response.Result().Header[key], should.Contain, value)
}

///////////////////////////////////////////////////////////////////////////////

type BadJSON struct{}

func (this *BadJSON) Error() string                { return "Implement the error interface." }
func (this *BadJSON) MarshalJSON() ([]byte, error) { return nil, errors.New("GOPHERS!") }

///////////////////////////////////////////////////////////////////////////////

func TestFirstNonBlank_WhenAllBlank_ReturnDefaultOfBlank(t *testing.T) {
	if actual := firstNonBlank("", "", ""); actual != "" {
		t.Error("Failed, expected '', got:", actual)
	}
}

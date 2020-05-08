package detour

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestResponsesFixture(t *testing.T) {
	gunit.Run(new(ResponsesFixture), t)
}

type ResponsesFixture struct {
	*gunit.Fixture
	body     string
	response *httptest.ResponseRecorder
	request  *http.Request
}

func (this *ResponsesFixture) Setup() {
	this.response = httptest.NewRecorder()
	this.request = httptest.NewRequest(http.MethodGet, "/?callback=hello", nil)
}
func (this *ResponsesFixture) render(renderers ...Renderer) {
	buffer := newResponseBuffer()
	renderer := CompoundRenderer(renderers)
	renderer.Render(buffer, this.request)
	buffer.flush(this.response)
	all, _ := ioutil.ReadAll(this.response.Result().Body)
	this.body = string(all)
}
func (this *ResponsesFixture) assertStatusCode(expected int) {
	this.So(this.response.Result().StatusCode, should.Equal, expected)
}
func (this *ResponsesFixture) assertStatusOK() {
	this.assertStatusCode(http.StatusOK)
}
func (this *ResponsesFixture) assertHeaders(pairs ...string) {
	headers := this.response.Result().Header
	for x := 0; x < len(pairs); x += 2 {
		key := pairs[x]
		value := pairs[x+1]
		this.So(headers[key], should.Contain, value)
	}
}
func (this *ResponsesFixture) assertNoResponseHeaders() {
	this.So(this.response.Result().Header, should.BeEmpty)
}
func (this *ResponsesFixture) assertBlankBody() {
	this.So(this.body, should.BeBlank)
}
func (this *ResponsesFixture) assertBody(expected string) {
	this.So(this.body, should.EqualTrimSpace, expected)
}
func (this *ResponsesFixture) assertBodyContains(substring string) {
	this.So(this.body, should.ContainSubstring, substring)
}
func (this *ResponsesFixture) assertUntouchedResponse() {
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBlankBody()
}

func (this *ResponsesFixture) TestNopRenderer() {
	this.render(NopRenderer{})
	this.assertUntouchedResponse()
}
func (this *ResponsesFixture) TestStatusCodeRenderer() {
	this.render(StatusCodeRenderer(http.StatusTeapot))
	this.assertStatusCode(http.StatusTeapot)
	this.assertNoResponseHeaders()
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestHeadersRenderer() {
	this.render(HeadersRenderer{"A": []string{"1"}, "B": []string{"2"}})
	this.assertStatusOK()
	this.assertHeaders(
		"A", "1",
		"B", "2",
	)
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestSetHeaderRenderer_OddLengthSlice_Panics() {
	this.So(func() { this.render(SetHeaderPairsRenderer{"x-smarty" /* missing value! */}) }, should.Panic)
}
func (this *ResponsesFixture) TestSetHeaderRenderer() {
	this.render(SetHeaderPairsRenderer{
		"x-smarty", "Streets",
		"x-hello", "world",
	})
	this.assertStatusOK()
	this.assertHeaders(
		"X-Smarty", "Streets",
		"X-Hello", "world",
	)
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestAddHeaderRenderer_OddLengthSlice_Panics() {
	this.So(func() { this.render(AddHeaderPairsRenderer{"x-smarty" /* missing value! */}) }, should.Panic)
}
func (this *ResponsesFixture) TestAddHeaderRenderer() {
	this.render(AddHeaderPairsRenderer{
		"x-smarty", "Streets",
		"x-smarty", "'s treats",
	})
	this.assertStatusOK()
	this.assertHeaders(
		"X-Smarty", "Streets",
		"X-Smarty", "'s treats",
	)
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestCookieRenderer() {
	cookie := &http.Cookie{
		Name:    "name",
		Value:   "value",
		Domain:  "domain",
		Expires: time.Now(),
	}
	this.render(CookieRenderer(*cookie))
	this.assertStatusOK()
	this.assertHeaders("Set-Cookie", cookie.String())
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestRedirect() {
	this.render(
		StatusCodeRenderer(http.StatusTemporaryRedirect),
		RedirectRenderer("https://smartystreets.com/redirect"),
	)
	this.assertStatusCode(http.StatusTemporaryRedirect)
	this.assertHeaders("Location", "https://smartystreets.com/redirect")
	this.assertBody(`<a href="https://smartystreets.com/redirect">Temporary Redirect</a>.`)
}
func (this *ResponsesFixture) TestBytesBodyRenderer() {
	this.render(BytesBodyRenderer("Hello, world!"))
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("Hello, world!")
}
func (this *ResponsesFixture) TestStringBodyRenderer() {
	this.render(StringBodyRenderer("Hello, world!"))
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("Hello, world!")
}
func (this *ResponsesFixture) TestReaderBodyRenderer() {
	this.render(ReaderBodyRenderer{Reader: strings.NewReader("Hello, world!")})
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("Hello, world!")
}
func (this *ResponsesFixture) TestDiagnosticBodyRenderer() {
	this.render(
		StatusCodeRenderer(http.StatusTeapot),
		DiagnosticBodyRenderer("Hello, world!"),
	)
	this.assertStatusCode(http.StatusTeapot)
	this.assertHeaders(
		"Content-Type", "text/plain; charset=utf-8",
		"X-Content-Type-Options", "nosniff",
	)
	this.assertBodyContains("418 Hello, world!")
	this.assertBodyContains("Raw Request:")
}
func (this *ResponsesFixture) TestXMLBodyRenderer() {
	this.render(XMLBodyRenderer{Content: "Hello, world!"})
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("<string>Hello, world!</string>")
}
func (this *ResponsesFixture) TestJSONBodyRenderer() {
	this.render(JSONBodyRenderer{Content: []int{1, 2, 3}})
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("[1,2,3]")
}
func (this *ResponsesFixture) TestJSONBodyRenderer_Indentation() {
	this.render(JSONBodyRenderer{Content: []int{1, 2, 3}, Indent: "  "})
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("[\n  1,\n  2,\n  3\n]")
}
func (this *ResponsesFixture) TestJSONBodyRenderer_JSONP() {
	this.render(JSONBodyRenderer{Content: []int{1, 2, 3}, JSONp: true})
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("hello([1,2,3])")
}
func (this *ResponsesFixture) TestJSONBodyRenderer_JSONP_Indent() {
	this.render(JSONBodyRenderer{Content: []int{1, 2, 3}, JSONp: true, Indent: "  "})
	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("hello([\n  1,\n  2,\n  3\n])")
}
func (this *ResponsesFixture) TestJSONBodyRenderer_JSONP_NoCallback() {
	query := this.request.URL.Query()
	query.Del("callback")
	this.request.URL.RawQuery = query.Encode()

	this.render(JSONBodyRenderer{Content: []int{1, 2, 3}, JSONp: true})

	this.assertStatusOK()
	this.assertNoResponseHeaders()
	this.assertBody("[1,2,3]")
}
func (this *ResponsesFixture) TestIfElseRendering_True() {
	this.render(IfElseRenderer(
		true,
		StatusCodeRenderer(http.StatusTeapot),
		StatusCodeRenderer(http.StatusInternalServerError)),
	)
	this.assertStatusCode(http.StatusTeapot)
	this.assertNoResponseHeaders()
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestIfElseRendering_False() {
	this.render(IfElseRenderer(
		false,
		StatusCodeRenderer(http.StatusTeapot),
		StatusCodeRenderer(http.StatusInternalServerError)),
	)
	this.assertStatusCode(http.StatusInternalServerError)
	this.assertNoResponseHeaders()
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestIfRendering_True() {
	this.render(IfRenderer(true, StatusCodeRenderer(http.StatusTeapot)))
	this.assertStatusCode(http.StatusTeapot)
	this.assertNoResponseHeaders()
	this.assertBlankBody()
}
func (this *ResponsesFixture) TestIfRendering_False() {
	this.render(IfRenderer(false, StatusCodeRenderer(http.StatusTeapot)))
	this.assertUntouchedResponse()
}

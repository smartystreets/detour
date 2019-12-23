package detour

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

/* ------------------------------------------------------------------------- */

type CompoundRenderer []Renderer

func (this CompoundRenderer) Render(response http.ResponseWriter, request *http.Request) {
	for _, renderer := range this {
		renderer.Render(response, request)
	}
}

/* ------------------------------------------------------------------------- */

func IfElseRenderer(condition bool, True, False Renderer) Renderer {
	if condition {
		return True
	} else {
		return False
	}
}

/* ------------------------------------------------------------------------- */

func IfRenderer(condition bool, renderer Renderer) Renderer {
	return IfElseRenderer(condition, renderer, NopRenderer{})
}

/* ------------------------------------------------------------------------- */

type NopRenderer struct{}

func (this NopRenderer) Render(http.ResponseWriter, *http.Request) {}

/* ------------------------------------------------------------------------- */

type StatusCodeRenderer int

func (this StatusCodeRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(int(this))
}

/* ------------------------------------------------------------------------- */

type HeadersRenderer http.Header

func (this HeadersRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	copyHeaders(http.Header(this), response.Header())
}

/* ------------------------------------------------------------------------- */

type SetHeaderPairsRenderer []string

func (this SetHeaderPairsRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	if len(this)%2 != 0 {
		panic("odd length")
	}
	header := response.Header()
	for x := 0; x < len(this); x += 2 {
		header.Set(this[x], this[x+1])
	}
}

/* ------------------------------------------------------------------------- */

type AddHeaderPairsRenderer []string

func (this AddHeaderPairsRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	if len(this)%2 != 0 {
		panic("odd length")
	}
	header := response.Header()
	for x := 0; x < len(this); x += 2 {
		header.Add(this[x], this[x+1])
	}
}

/* ------------------------------------------------------------------------- */

type CookieRenderer http.Cookie

func (this CookieRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	cookie := http.Cookie(this)
	http.SetCookie(response, &cookie)
}

/* ------------------------------------------------------------------------- */

type RedirectRenderer string

func (this RedirectRenderer) Render(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, string(this), response.(*responseBuffer).StatusCode())
}

/* ------------------------------------------------------------------------- */

type BytesBodyRenderer []byte

func (this BytesBodyRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	_, _ = response.Write(this)
}

/* ------------------------------------------------------------------------- */

type StringBodyRenderer string

func (this StringBodyRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(response, string(this))
}

/* ------------------------------------------------------------------------- */

type ReaderBodyRenderer struct{ io.Reader }

func (this ReaderBodyRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	_, _ = io.Copy(response, this.Reader)
}

/* ------------------------------------------------------------------------- */

type DiagnosticBodyRenderer string

func (this DiagnosticBodyRenderer) Render(response http.ResponseWriter, request *http.Request) {
	statusCode := response.(*responseBuffer).StatusCode()
	dump, _ := httputil.DumpRequest(request, false)
	requestDump := formatRequestDump(string(dump))
	message := fmt.Sprintf(diagnosticTemplate, statusCode, this, requestDump, disclaimer)
	http.Error(response, message, statusCode)
}

/* ------------------------------------------------------------------------- */

type XMLBodyRenderer struct{ Content interface{} }

func (this XMLBodyRenderer) Render(response http.ResponseWriter, _ *http.Request) {
	_ = xml.NewEncoder(response).Encode(this.Content)
}

/* ------------------------------------------------------------------------- */

type JSONBodyRenderer struct {
	Content interface{}
	Indent  string
	JSONp   bool
}

func (this JSONBodyRenderer) Render(response http.ResponseWriter, request *http.Request) {
	if this.JSONp {
		_, _ = io.WriteString(response, request.URL.Query().Get("callback"))
		_, _ = io.WriteString(response, "(")
	}

	encoder := json.NewEncoder(response)
	if this.Indent != "" {
		encoder.SetIndent("", this.Indent)
	}
	_ = encoder.Encode(this.Content)

	if this.JSONp {
		_, _ = io.WriteString(response, ")")
	}
}

/* ------------------------------------------------------------------------- */

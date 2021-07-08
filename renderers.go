package detour

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
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
	callback := request.URL.Query().Get("callback")
	canRenderJSONp := this.JSONp && callback != ""

	if canRenderJSONp {
		_, _ = io.WriteString(response, callback)
		_, _ = io.WriteString(response, "(")
	}

	_, _ = response.Write(this.renderContent())

	if canRenderJSONp {
		_, _ = io.WriteString(response, ")")
	}
}
func (this JSONBodyRenderer) renderContent() (content []byte) {
	if this.Indent != "" {
		content, _ = json.MarshalIndent(this.Content, "", this.Indent)
	} else {
		content, _ = json.Marshal(this.Content)
	}
	return content
}

/* ------------------------------------------------------------------------- */

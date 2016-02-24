package detour

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type (
	StatusCodeResult struct {
		StatusCode int
		Message    string
	}
	ContentResult struct {
		StatusCode  int
		ContentType string
		Content     string
	}
	DiagnosticResult struct {
		StatusCode int
		Message    string
	}
	BinaryResult struct {
		StatusCode  int
		ContentType string
		Content     []byte
	}
	JSONResult struct {
		StatusCode  int
		ContentType string
		Content     interface{}
	}
	JSONPResult struct {
		StatusCode    int
		ContentType   string
		Content       interface{}
	}
	ValidationResult struct {
		Failure1 error
		Failure2 error
		Failure3 error
		Failure4 error
	}
	ErrorResult struct {
		StatusCode int
		Error1     error
		Error2     error
		Error3     error
		Error4     error
	}
	CookieResult struct {
		Cookie1 *http.Cookie
		Cookie2 *http.Cookie
		Cookie3 *http.Cookie
		Cookie4 *http.Cookie
	}
)

func (this *StatusCodeResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentTypeAndStatusCode(response, this.StatusCode, plaintextContentType)
	response.Write([]byte(this.Message))
}

func (this *ContentResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, plaintextContentType)
	writeContentTypeAndStatusCode(response, this.StatusCode, contentType)
	response.Write([]byte(this.Content))
}

func (this *DiagnosticResult) Render(response http.ResponseWriter, request *http.Request) {
	message := composeErrorBody(request, this.Message, this.StatusCode)
	http.Error(response, message, this.StatusCode)
}

func (this *BinaryResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, octetStreamContentType)
	writeContentTypeAndStatusCode(response, this.StatusCode, contentType)
	response.Write(this.Content)
}

func (this *JSONResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, jsonContentType)
	writeContentType(response, contentType)
	serializeAndWrite(response, this.StatusCode, this.Content)
}
func (this *JSONPResult) Render(response http.ResponseWriter, request *http.Request) {
	contentType := selectContentType(this.ContentType, jsonContentType)
	writeContentType(response, contentType)
	callbackLabel := request.URL.Query().Get("callback") // We don't call request.ParseForm in every case so using the URL.Query() is safer.
	serializeAndWriteJSONP(response, this.StatusCode, this.Content, callbackLabel)
}

func (this *ValidationResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, jsonContentType)

	var failures Errors
	failures = failures.Append(this.Failure1)
	failures = failures.Append(this.Failure2)
	failures = failures.Append(this.Failure3)
	failures = failures.Append(this.Failure4)

	serializeAndWrite(response, 422, failures)
}

func (this *ErrorResult) Render(response http.ResponseWriter, request *http.Request) {
	writeContentType(response, jsonContentType)

	var failures Errors
	failures = failures.Append(this.Error1)
	failures = failures.Append(this.Error2)
	failures = failures.Append(this.Error3)
	failures = failures.Append(this.Error4)

	serializeAndWrite(response, this.StatusCode, failures)
}

func (this *CookieResult) Render(response http.ResponseWriter, request *http.Request) {
	for _, cookie := range []*http.Cookie{this.Cookie1, this.Cookie2, this.Cookie3, this.Cookie4} {
		if cookie != nil {
			http.SetCookie(response, cookie)
		}
	}
}

func selectContentType(values ...string) string {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}

	return ""
}

func writeContentTypeAndStatusCode(response http.ResponseWriter, statusCode int, contentType string) {
	writeContentType(response, contentType)
	response.WriteHeader(defaultToHTTPStatusOK(statusCode))
}
func writeContentType(response http.ResponseWriter, contentType string) {
	if len(contentType) > 0 {
		response.Header().Set(contentTypeHeader, contentType) // doesn't get written unless status code is written last!
	}
}

func serializeAndWrite(response http.ResponseWriter, statusCode int, content interface{}) {
	if content, err := json.Marshal(content); err == nil {
		writeContent(response, statusCode, content)
	} else {
		writeError(response)
	}
}
func serializeAndWriteJSONP(response http.ResponseWriter, statusCode int, content interface{}, label string) {
	if len(label) == 0 {
		serializeAndWrite(response, statusCode, content)
	} else if content, err := json.Marshal(content); err == nil {
		buffer := bytes.NewBufferString(label)
		buffer.WriteString("(")
		buffer.Write(content)
		buffer.WriteString(")")
		writeContent(response, statusCode, buffer.Bytes())
	} else {
		writeError(response)
	}
}
func writeContent(response http.ResponseWriter, statusCode int, content []byte) {
	response.WriteHeader(defaultToHTTPStatusOK(statusCode))
	response.Write(content)
}
func writeError(response http.ResponseWriter) {
	response.WriteHeader(http.StatusInternalServerError)
	errContent := make(Errors, 0).Append(SimpleInputError("Marshal failure", "HTTP Response"))
	content, _ := json.Marshal(errContent)
	response.Write(content)
}

func defaultToHTTPStatusOK(statusCode int) int {
	if statusCode == 0 {
		return http.StatusOK
	}
	return statusCode
}

const (
	contentTypeHeader      = "Content-Type"
	jsonContentType        = "application/json; charset=utf-8"
	octetStreamContentType = "application/octet-stream"
	plaintextContentType   = "text/plain; charset=utf-8"
)

///////////////////////////////////////////////////////////////////////////////

// Sources for ascii art in disclaimer:
// - http://patorjk.com/software/taag/#p=display&f=Standard&t=SmartyStreets
// - http://www.chris.com/ascii/index.php?art=art%20and%20design/borders
//
// It's ok to ignore the error here. A blank disclaimer in that case isn't a show-stopper.
var disclaimer, _ = base64.StdEncoding.DecodeString("CiAgLi0tLiAgICAgIC4tJy4gICAgICAuLS0uICAgICAgLi0tLiAgICAgIC4tLS4gICAgICAuLS0uICAgICAgLmAtLiAgICAgCjo6Ojo6Llw6Ojo6Ojo6Oi5cOjo6Ojo6OjouXDo6Ojo6Ojo6Llw6Ojo6Ojo6Oi5cOjo6Ojo6OjouXDo6Ojo6Ojo6Llw6Ojo6CicgICAgICBgLS0nICAgICAgYC4tJyAgICAgIGAtLScgICAgICBgLS0nICAgICAgYC0tJyAgICAgIGAtLicgICAgICBgLS0nCgoKICAuLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0uCiAvICAuLS4gICAgICAgICAgICAgICAgLS0tLSBESVNDTEFJTUVSIC0tLS0gICAgICAgICAgICAgICAgICAgICAuLS4gIFwKfCAgLyAgIFwgICAgIFRoZSBvdXRwdXQgeW91IHNlZSBoZXJlIGhhcyBiZWVuIGdlbmVyYXRlZCBhcyAgICAgLyAgIFwgIHwKfCB8XF8uICB8IGNvbnZlbmllbmNlIHRvIGFpZCBpbiBkZWJ1Z2dpbmcgY2xpZW50IGFwcGxpY2F0aW9ucy58ICAgIC98IHwKfFx8ICB8IC98ICBJdCBpcyBzdWJqZWN0IHRvIGNoYW5nZSB3aXRob3V0IG5vdGljZSBhbmQgZm9yIGFueSB8XCAgfCB8L3wKfCBgLS0tJyB8ICByZWFzb24uIFBsZWFzZSBwcm9ncmFtIHlvdXIgYXBwbGljYXRpb25zIHRvIGNoZWNrICB8IGAtLS0nIHwKfCAgICAgICB8ICAgdGhlIEhUVFAgU3RhdHVzIENvZGUgYW5kIG9ubHkgcGFyc2UgdGhlIGJvZHkgaW4gICB8ICAgICAgIHwKfCAgICAgICB8ICAgICBjYXNlIG9mIGEgc3RhdHVzIGNvZGUgb2YgJzIwMCcuIFRoYXQgaXMgYWxsLiAgICB8ICAgICAgIHwKfCAgICAgICB8LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS18ICAgICAgIHwKXCAgICAgICB8ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB8ICAgICAgIC8KIFwgICAgIC8gICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgXCAgICAgLwogIGAtLS0nICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgYC0tLScKCgogIC4tLS4gICAgICAuLScuICAgICAgLi0tLiAgICAgIC4tLS4gICAgICAuLS0uICAgICAgLi0tLiAgICAgIC5gLS4gICAgIAo6Ojo6Oi5cOjo6Ojo6OjouXDo6Ojo6Ojo6Llw6Ojo6Ojo6Oi5cOjo6Ojo6OjouXDo6Ojo6Ojo6Llw6Ojo6Ojo6Oi5cOjo6OgonICAgICAgYC0tJyAgICAgIGAuLScgICAgICBgLS0nICAgICAgYC0tJyAgICAgIGAtLScgICAgICBgLS4nICAgICAgYC0tJwoKCiAgX19fXyAgICAgICAgICAgICAgICAgICAgICAgXyAgICAgICAgIF9fX18gIF8gICAgICAgICAgICAgICAgIF8gICAgICAgCiAvIF9fX3wgXyBfXyBfX18gICBfXyBfIF8gX198IHxfIF8gICBfLyBfX198fCB8XyBfIF9fIF9fXyAgX19ffCB8XyBfX18gCiBcX19fIFx8ICdfIGAgXyBcIC8gX2AgfCAnX198IF9ffCB8IHwgXF9fXyBcfCBfX3wgJ19fLyBfIFwvIF8gXCBfXy8gX198CiAgX19fKSB8IHwgfCB8IHwgfCAoX3wgfCB8ICB8IHxffCB8X3wgfF9fXykgfCB8X3wgfCB8ICBfXy8gIF9fLyB8X1xfXyBcCiB8X19fXy98X3wgfF98IHxffFxfXyxffF98ICAgXF9ffFxfXywgfF9fX18vIFxfX3xffCAgXF9fX3xcX19ffFxfX3xfX18vCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHxfX18vICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgIAogICAgICAuLi4uICAgICAgICAgICAuLi4uICAgICAgICAgICAuLi4uICAgICAgICAgICAuLi4uICAgICAgICAgICAuLi4uCiAgICAgfHwgICAgICAgICAgICAgfHwgICAgICAgICAgICAgfHwgICAgICAgICAgICAgfHwgICAgICAgICAgICAgfHwKIC8iIiJsfFwgICAgICAgIC8iIiJsfFwgICAgICAgIC8iIiJsfFwgICAgICAgIC8iIiJsfFwgICAgICAgIC8iIiJsfFwKL19fX19fX19cICAgICAgL19fX19fX19cICAgICAgL19fX19fX19cICAgICAgL19fX19fX19cICAgICAgL19fX19fX19cCnwgIC4tLiAgfC0tLS0tLXwgIC4tLiAgfC0tLS0tLXwgIC4tLiAgfC0tLS0tLXwgIC4tLiAgfC0tLS0tLXwgIC4tLiAgfC0tLS0tLQogX198THxfX3wgLi0tLiAgX198THxfX3wgLi0tLiB8X198THxfX3wgLi0tLiB8X198THxfX3wgLi0tLiB8X198THxfX3wgLi0tLgpfXCAgXFxwX19gby1vJ19fXCAgXFxwX19gby1vJ19fXCAgXFxwX19gby1vJ19fXCAgXFxwX19gby1vJ19fXCAgXFxwX19gby1vJ18KLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tCi0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLQo=")

func composeErrorBody(request *http.Request, message string, status int) string {
	dump, _ := httputil.DumpRequest(request, false)
	return fmt.Sprintf("%d %s\n\nRaw Request:\n\n%s\n\n%s", status, message, string(dump), disclaimer)
}

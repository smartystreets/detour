package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/smartystreets/detour"
)

func main() {
	controller := &Controller{} // load initial state from environment, or data store...

	// You can use any routing mechanism you want here. detour.New(..)
	// gives back http.Handler. Pass references to controller methods
	// to detour.New(..) and all the Bind(), Sanitize(), Validate(),
	// and Renderer function will be called on the input parameter of the
	// controller function.
	http.Handle("/hello", detour.New(controller.SayHello))
	// http.Handle("/world", detour.New(controller.SomeOtherAction))

	address := "127.0.0.1:8080"
	log.Println("Listening on", address)
	http.ListenAndServe(address, nil)

	/*
		$ go run main.go

		<switch to another terminal window>

		$ curl -v "http://localhost:8080/hello"
		*   Trying ::1...
		* Connected to localhost (::1) port 8080 (#0)
		> GET /hello HTTP/1.1
		> Host: localhost:8080
		> User-Agent: curl/7.43.0
		>
		< HTTP/1.1 422 status code 422
		< Content-Type: application/json; charset=utf-8
		< Date: Wed, 25 Nov 2015 23:37:44 GMT
		< Content-Length: 55
		<
		* Connection #0 to host localhost left intact
		[{"fields":["name"],"message":"The field is required"}]

		$ curl -v "http://localhost:8080/hello?name=mike"
		*   Trying ::1...
		* Connected to localhost (::1) port 8080 (#0)
		> GET /hello?name=Mike HTTP/1.1
		> Host: localhost:8080
		> User-Agent: curl/7.43.0
		>
		< HTTP/1.1 202 Accepted
		< Content-Type: text/plain
		< Date: Wed, 25 Nov 2015 23:38:17 GMT
		< Content-Length: 12
		<
		* Connection #0 to host localhost left intact
		Hello, Mike!
	*/
}

///////////////////////////////////////////////////////////////////////////////

type Controller struct {
}

// SayHello is a controller action that, when called by the ServeHTTP method of
// the library's actionHandler will receive a SalutationInputModel. By this time the
// input's Bind(), Sanitize(), Validate(), and Error() methods will have already
// been called. The returned detour.Renderer will be written to the actual
// http.ResponseWriter by the ServeHTTP method of the actionHandler. The detour
// package provides various types that implement the Renderer interface. Users of
// this package may also supply their own types that implement the Renderer interface.
func (this *Controller) SayHello(input *SalutationInputModel) detour.Renderer {
	// This ContentResult will be serialized to the http.ResponseWriter in actionHandler.
	return detour.ContentResult{
		StatusCode:  http.StatusAccepted,
		ContentType: "text/plain",
		Content:     fmt.Sprintf("Hello, %s!", input.Name),
	}
}

///////////////////////////////////////////////////////////////////////////////

type SalutationInputModel struct {
	Name string
}

// Bind receives the actual *http.Request and pulls off the bits and pieces
// necessary to populate the fields on this. This could be where you deserialize
// a JSON payload in the request.Body or access the request form, the query
// string, the URL, or whatever! Any error returned here will become
// an HTTP 400 (Bad Request) and will prevent the Validate method from being
// called along with any controller actions that expect/receive this type.
// Implementing this method is optional (but the only way to get information
// from the request onto the input model).
func (this *SalutationInputModel) Bind(request *http.Request) error {
	// request.ParseForm() will have already been called in actionHandler.
	this.Name = request.Form.Get("name")
	return nil
}

// Sanitize performs, post processing (cleanup) on the data bound from the *http.Request.
// The reason for splitting apart Bind and Sanitize is to allow the Sanitize logic to be
// tested independently of the *http.Request which is received during Bind.
// Sanitize returns no error, but could save errors for Validate to return if needed.
// Implementing this method is optional.
func (this *SalutationInputModel) Sanitize() {
	this.Name = strings.TrimSpace(strings.Title(this.Name))
}

// Validate inspects the fields populated by Bind() (and cleaned by Sanitize) to ensure
// that they are semantically correct. Any error returned from this function will result
// in an HTTP 422 (Unprocessable Entity) and will skip any controller actions that
// expect/receive this type.
// Implementing this method is optional.
func (this *SalutationInputModel) Validate() error {
	var errors detour.Errors
	if len(this.Name) == 0 {
		errors = errors.Append(detour.SimpleInputError("The field is required", "name"))
	}
	return errors
}

// Error allows an HTTP 500 Internal Server Error to be returned if conditions are right for that.
// Implementing this method is optional. If implemented, a true value causes the HTTP 500 and
// short-circuits any calls to controller actions that receive this type.
func (this *SalutationInputModel) Error() bool {
	return false
}

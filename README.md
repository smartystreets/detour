[![Build Status](https://travis-ci.org/smartystreets/detour.svg?branch=master)](https://travis-ci.org/smartystreets/detour)

# detour
--
    import "."

package detour offers an alternate, MVC-based, approach to HTTP applications.
Rather than writing traditional http.Handlers you define input models that have
optional Bind(), Sanitize(), and Validate() methods and which can be passed into
methods on structs which return a Renderer. Each of these concepts is glued
together by the library's actionHandler struct via the New() function. See the
example folder for a complete example. Requires Go 1.7+

## Usage

#### func  Bind

```go
func Bind(request *http.Request, message interface{}) error
```
Bind is exported for use in testing.

#### func  CompoundInputError

```go
func CompoundInputError(message string, fields ...string) error
```

#### func  New

```go
func New(controllerAction interface{}) http.Handler
```

#### func  NewFromFactory

```go
func NewFromFactory(inputModelFactory createModel, controllerAction interface{}) http.Handler
```

#### func  SimpleInputError

```go
func SimpleInputError(message, field string) error
```

#### type BinaryResult

```go
type BinaryResult struct {
	StatusCode  int
	ContentType string
	Content     []byte
}
```


#### func (BinaryResult) Render

```go
func (this BinaryResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type BindJSON

```go
type BindJSON interface {
	BindJSON() bool
}
```


#### type Binder

```go
type Binder interface {
	Bind(*http.Request) error
}
```


#### type ContentResult

```go
type ContentResult struct {
	StatusCode  int
	ContentType string
	Content     string
	Headers     map[string]string // TODO: do we even need/use this?
}
```


#### func (ContentResult) Render

```go
func (this ContentResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type CookieResult

```go
type CookieResult struct {
	Cookie1 *http.Cookie
	Cookie2 *http.Cookie
	Cookie3 *http.Cookie
	Cookie4 *http.Cookie
}
```


#### func (CookieResult) Render

```go
func (this CookieResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type DiagnosticError

```go
type DiagnosticError struct {
	HTTPStatusCode int
}
```


#### func  NewDiagnosticError

```go
func NewDiagnosticError(message string) *DiagnosticError
```

#### func (*DiagnosticError) Error

```go
func (this *DiagnosticError) Error() string
```

#### func (*DiagnosticError) StatusCode

```go
func (this *DiagnosticError) StatusCode() int
```

#### type DiagnosticResult

```go
type DiagnosticResult struct {
	StatusCode int
	Message    string
}
```


#### func (DiagnosticResult) Render

```go
func (this DiagnosticResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type ErrorCode

```go
type ErrorCode interface {
	error
	StatusCode() int
}
```


#### type ErrorResult

```go
type ErrorResult struct {
	StatusCode int
	Error1     error
	Error2     error
	Error3     error
	Error4     error
}
```


#### func (ErrorResult) Render

```go
func (this ErrorResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type Errors

```go
type Errors []error
```


#### func (Errors) Append

```go
func (this Errors) Append(err error) Errors
```

#### func (Errors) AppendIf

```go
func (this Errors) AppendIf(err error, condition bool) Errors
```

#### func (Errors) Error

```go
func (this Errors) Error() string
```

#### func (Errors) MarshalJSON

```go
func (this Errors) MarshalJSON() ([]byte, error)
```

#### func (Errors) StatusCode

```go
func (this Errors) StatusCode() int
```

#### type InputError

```go
type InputError struct {
	Fields         []string `json:"fields"`
	Message        string   `json:"message"`
	HTTPStatusCode int      `json:"-"`
}
```


#### func (*InputError) Error

```go
func (this *InputError) Error() string
```

#### func (*InputError) StatusCode

```go
func (this *InputError) StatusCode() int
```

#### type JSONPResult

```go
type JSONPResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
	Indent      string
}
```


#### func (JSONPResult) Render

```go
func (this JSONPResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type JSONResult

```go
type JSONResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
	Indent      string
}
```


#### func (JSONResult) Render

```go
func (this JSONResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type RedirectResult

```go
type RedirectResult struct {
	Location   string
	StatusCode int
}
```


#### func (RedirectResult) Render

```go
func (this RedirectResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type Renderer

```go
type Renderer interface {
	Render(http.ResponseWriter, *http.Request)
}
```


#### type Sanitizer

```go
type Sanitizer interface {
	Sanitize()
}
```


#### type ServerError

```go
type ServerError interface {
	Error() bool
}
```


#### type StatusCodeResult

```go
type StatusCodeResult struct {
	StatusCode int
	Message    string
}
```


#### func (StatusCodeResult) Render

```go
func (this StatusCodeResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type ValidationResult

```go
type ValidationResult struct {
	Failure1 error
	Failure2 error
	Failure3 error
	Failure4 error
}
```


#### func (ValidationResult) Render

```go
func (this ValidationResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type Validator

```go
type Validator interface {
	Validate() error
}
```

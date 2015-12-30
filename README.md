# detour
--
    import "github.com/smartystreets/detour"

package detour offers an alternate, MVC-based, approach to HTTP applications.
Rather than writing traditional http.Handlers you define input models that have
optional Bind(), Sanitize(), and Validate() methods and which can be passed into
methods on structs which return a Renderer. Each of these concepts is glued
together by the library's ActionHandler struct via the New() function. See the
example folder for a complete example.

## Usage

#### func  CompoundInputError

```go
func CompoundInputError(message string, fields ...string) error
```

#### func  New

```go
func New(controllerAction interface{}) http.Handler
```

#### func  SimpleInputError

```go
func SimpleInputError(message, field string) error
```

#### type ActionHandler

```go
type ActionHandler struct {
}
```


#### func (*ActionHandler) ServeHTTP

```go
func (this *ActionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request)
```

#### type BinaryResult

```go
type BinaryResult struct {
	StatusCode  int
	ContentType string
	Content     []byte
}
```


#### func (*BinaryResult) Render

```go
func (this *BinaryResult) Render(response http.ResponseWriter, request *http.Request)
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
}
```


#### func (*ContentResult) Render

```go
func (this *ContentResult) Render(response http.ResponseWriter, request *http.Request)
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


#### func (*CookieResult) Render

```go
func (this *CookieResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type CreateModel

```go
type CreateModel func() interface{}
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


#### func (*ErrorResult) Render

```go
func (this *ErrorResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type Errors

```go
type Errors []error
```


#### func (Errors) Append

```go
func (this Errors) Append(err error) Errors
```

#### func (Errors) Error

```go
func (this Errors) Error() string
```

#### type InputError

```go
type InputError struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}
```


#### func (*InputError) Error

```go
func (this *InputError) Error() string
```

#### type JSONResult

```go
type JSONResult struct {
	StatusCode  int
	ContentType string
	Content     interface{}
}
```


#### func (*JSONResult) Render

```go
func (this *JSONResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type MonadicAction

```go
type MonadicAction func(interface{}) Renderer
```


#### type NiladicAction

```go
type NiladicAction func() Renderer
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


#### func (*StatusCodeResult) Render

```go
func (this *StatusCodeResult) Render(response http.ResponseWriter, request *http.Request)
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


#### func (*ValidationResult) Render

```go
func (this *ValidationResult) Render(response http.ResponseWriter, request *http.Request)
```

#### type Validator

```go
type Validator interface {
	Validate() error
}
```

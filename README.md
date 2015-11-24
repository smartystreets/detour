# binding
--
    import "github.com/smartystreets/binding"


## Usage

#### func  ComplexValidationError

```go
func ComplexValidationError(message string, fields ...string) error
```

#### func  New

```go
func New(controllerAction interface{}) http.Handler
```

#### func  SimpleValidationError

```go
func SimpleValidationError(message, field string) error
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

#### type ValidationError

```go
type ValidationError struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}
```


#### func (*ValidationError) Error

```go
func (this *ValidationError) Error() string
```

#### type ValidationErrors

```go
type ValidationErrors []error
```


#### func (ValidationErrors) Append

```go
func (this ValidationErrors) Append(err error) ValidationErrors
```

#### func (ValidationErrors) Error

```go
func (this ValidationErrors) Error() string
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

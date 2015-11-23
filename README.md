# binding
--
    import "github.com/smartystreets/binding"


## Usage

#### func  ComplexValidationError

```go
func ComplexValidationError(message string, fields ...string) error
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


#### func  Typed

```go
func Typed(controllerAction interface{}) *ActionHandler
```

#### func  TypedFactory

```go
func TypedFactory(controllerAction interface{}, input InputFactory) *ActionHandler
```

#### func (*ActionHandler) ServeHTTP

```go
func (this *ActionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request)
```

#### type Binder

```go
type Binder interface {
	Bind(*http.Request) error
}
```


#### type ControllerAction

```go
type ControllerAction func(interface{}) Renderer
```


#### type InputFactory

```go
type InputFactory func() interface{}
```


#### type Renderer

```go
type Renderer interface {
	Render(http.ResponseWriter, *http.Request)
}
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

#### type Validator

```go
type Validator interface {
	Validate() error
}
```

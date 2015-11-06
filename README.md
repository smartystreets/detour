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

#### type Binder

```go
type Binder interface {
	Bind(request *http.Request) error
}
```


#### type DomainHandler

```go
type DomainHandler func(interface{}) http.Handler
```


#### type Handler

```go
type Handler interface {
	Handle(response http.ResponseWriter, request *http.Request, message interface{})
}
```


#### type InputModelFactory

```go
type InputModelFactory func() interface{}
```


#### type ModelBinder

```go
type ModelBinder struct {
}
```


#### func  NewDomainModelBinder

```go
func NewDomainModelBinder(input InputModelFactory, domain DomainHandler) *ModelBinder
```

#### func  NewModelBinderHandler

```go
func NewModelBinderHandler(input InputModelFactory, handler Handler) *ModelBinder
```

#### func (*ModelBinder) ServeHTTP

```go
func (this *ModelBinder) ServeHTTP(response http.ResponseWriter, request *http.Request)
```

#### type Translator

```go
type Translator interface {
	Translate() interface{}
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

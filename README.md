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
	Bind(*http.Request) error
}
```


#### type ControllerAction

```go
type ControllerAction func(http.ResponseWriter, *http.Request, interface{})
```


#### type DomainAction

```go
type DomainAction func(interface{}) http.Handler
```


#### type InputFactory

```go
type InputFactory func() interface{}
```


#### type ModelBinder

```go
type ModelBinder struct {
}
```


#### func  Domain

```go
func Domain(callback DomainAction, input InputFactory) *ModelBinder
```

#### func  Generic

```go
func Generic(callback ControllerAction, message interface{}) *ModelBinder
```

#### func  GenericFactory

```go
func GenericFactory(callback ControllerAction, input InputFactory) *ModelBinder
```

#### func  Typed

```go
func Typed(controllerAction interface{}) *ModelBinder
```

#### func  TypedFactory

```go
func TypedFactory(controllerAction interface{}, input InputFactory) *ModelBinder
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

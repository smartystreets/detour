package detour

import "encoding/json"

type Errors []error

func (this Errors) AppendIf(err error, condition bool) Errors {
	if condition {
		return this.Append(err)
	} else {
		return this
	}
}

func (this Errors) Append(err error) Errors {
	if err != nil {
		this = append(this, err)
	}
	return this
}

func (this Errors) Error() string {
	raw, _ := this.MarshalJSON()
	return string(raw)
}

func (this Errors) MarshalJSON() ([]byte, error) {
	var filtered []error
	for _, err := range this {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	return json.Marshal(filtered)
}

func (this Errors) StatusCode() int {
	for _, err := range this {
		if code, ok := err.(ErrorCode); ok {
			statusCode := code.StatusCode()
			if statusCode != 0 {
				return statusCode
			}
		}
	}
	return 0
}

type InputError struct {
	Fields         []string `json:"fields"`
	Message        string   `json:"message"`
	HTTPStatusCode int      `json:"-"`
}

func SimpleInputError(message, field string) error {
	return &InputError{Fields: []string{field}, Message: message}
}
func CompoundInputError(message string, fields ...string) error {
	return &InputError{Fields: fields, Message: message}
}
func (this *InputError) Error() string {
	raw, _ := json.Marshal(this)
	return string(raw)
}

func (this *InputError) StatusCode() int {
	return this.HTTPStatusCode
}

///////////////////////////////////////////////////////////////////////////////

type DiagnosticError struct {
	message        string
	HTTPStatusCode int
}

func NewDiagnosticError(message string) *DiagnosticError {
	return &DiagnosticError{message: message}
}

func (this *DiagnosticError) Error() string {
	return this.message
}

func (this *DiagnosticError) StatusCode() int {
	return this.HTTPStatusCode
}

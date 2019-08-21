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

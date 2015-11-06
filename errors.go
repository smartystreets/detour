package binding

import "encoding/json"

type ValidationErrors []error

func (this ValidationErrors) Append(err error) ValidationErrors {
	if err != nil {
		this = append(this, err)
	}
	return this
}

func (this ValidationErrors) Error() string {
	filtered := []error{}
	for _, err := range this {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	raw, _ := json.Marshal(filtered)
	return string(raw)
}

type ValidationError struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

func SimpleValidationError(message, field string) error {
	return &ValidationError{Fields: []string{field}, Message: message}
}
func ComplexValidationError(message string, fields ...string) error {
	return &ValidationError{Fields: fields, Message: message}
}
func (this *ValidationError) Error() string {
	return this.Message
}

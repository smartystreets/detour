package binding

import "encoding/json"

type Errors []error

func (this Errors) Append(err error) Errors {
	if err != nil {
		this = append(this, err)
	}
	return this
}

func (this Errors) Error() string {
	filtered := []error{}
	for _, err := range this {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	raw, _ := json.Marshal(filtered)
	return string(raw)
}

type InputError struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

func SimpleInputError(message, field string) error {
	return &InputError{Fields: []string{field}, Message: message}
}
func CompoundInputError(message string, fields ...string) error {
	return &InputError{Fields: fields, Message: message}
}
func (this *InputError) Error() string {
	return this.Message
}

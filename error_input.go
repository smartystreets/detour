package detour

import "encoding/json"

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

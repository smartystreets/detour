package detour

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

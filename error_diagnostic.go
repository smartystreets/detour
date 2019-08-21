package detour

// Deprecated in favor of DiagnosticErrors
type DiagnosticError struct {
	message        string
	HTTPStatusCode int
}

// Deprecated
func NewDiagnosticError(message string) *DiagnosticError {
	return &DiagnosticError{message: message}
}

// Deprecated
func (this *DiagnosticError) Error() string {
	return this.message
}

// Deprecated
func (this *DiagnosticError) StatusCode() int {
	return this.HTTPStatusCode
}

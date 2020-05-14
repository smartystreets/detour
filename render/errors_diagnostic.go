package render

import (
	"fmt"
	"strconv"
	"strings"
)

type DiagnosticErrors []error

func (this DiagnosticErrors) Append(err error) DiagnosticErrors {
	if err != nil {
		this = append(this, err)
	}
	return this
}

func (this DiagnosticErrors) AppendIf(err error, condition bool) DiagnosticErrors {
	if condition {
		return this.Append(err)
	}
	return this
}

func (this DiagnosticErrors) Error() string {
	return fmt.Sprintf("Errors:\n\n") + this.list()
}

func (this DiagnosticErrors) list() string {
	var builder strings.Builder
	for e, err := range this {
		if len(this) == 1 {
			builder.WriteString("- ")
		} else {
			builder.WriteString(strconv.Itoa(e + 1))
			builder.WriteString(". ")
		}
		builder.WriteString(err.Error())
		builder.WriteString("\n")
	}
	return strings.TrimSpace(builder.String())
}

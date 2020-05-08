package detour

import (
	"context"
	"net/http"
)

type (
	Binder interface {
		Bind(*http.Request) error
	}

	BindContext interface {
		BindContext(context.Context)
	}

	BindJSON interface {
		BindJSON() bool
	}

	Sanitizer interface {
		Sanitize()
	}

	Validator interface {
		Validate() error
	}

	ServerError interface {
		Error() bool
	}

	Renderer interface {
		Render(http.ResponseWriter, *http.Request)
	}

	ErrorCode interface {
		error
		StatusCode() int
	}
)

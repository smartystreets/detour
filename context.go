package detour

import "context"

type ContextBinder struct {
	Context context.Context
}

func (this *ContextBinder) BindContext(context context.Context) {
	this.Context = context
}

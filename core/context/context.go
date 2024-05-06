package context

import "context"

type IContext interface {
	Context() context.Context
}

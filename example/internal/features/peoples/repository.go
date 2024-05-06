package peoples

import "context"

type IRepository interface {
	Cache(ctx context.Context, key string) (string, error)
}

package peoples

import "context"

type IService interface {
	Cache(ctx context.Context, key string) (string, error)
}

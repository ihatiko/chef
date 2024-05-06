package planets

import (
	"context"
	"example/internal/features/planets/types"
)

type IReadRepository interface {
	Get(ctx context.Context) ([]types.Planet, error)
}

type IWriteRepository interface {
	Update() error
}

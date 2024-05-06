package read

import (
	"context"
	"example/internal/features/planets/repository/read/queries"
	"example/internal/features/planets/types"
)

func (r repositoty) Get(ctx context.Context) ([]types.Planet, error) {
	result := []types.Planet{}
	return result, r.client.Db.GetContext(ctx, &result, queries.GetAllPlanets)
}

package read

import (
	"context"
	"example/internal/features/planets/repository/read/queries"
	"example/internal/features/planets/types"
)

func (r repository) Get(ctx context.Context) ([]types.Planet, error) {
	var result []types.Planet
	return result, r.client.Db.GetContext(ctx, &result, queries.GetAllPlanets)
}

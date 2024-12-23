package read

import (
	"example/internal/features/planets"

	"github.com/ihatiko/chef/components/clients/postgresql"
)

type repository struct {
	client postgresql.Client
}

func New(db postgresql.Client) planets.IReadRepository {
	return &repository{db}
}

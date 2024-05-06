package read

import (
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/temple/infrastucture/postgresql"
)

type repositoty struct {
	client postgresql.Client
}

func (r repositoty) Get() error {
	//TODO implement me
	panic("implement me")
}

func New(db postgresql.Client) planets.IReadRepository {
	r := new(repositoty)
	r.client = db
	return r
}

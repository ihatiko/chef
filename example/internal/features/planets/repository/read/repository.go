package read

import (
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/temple/infrastucture/postgresql"
)

type Repositoty struct {
	client postgresql.Client
}

func (r Repositoty) Get() error {
	//TODO implement me
	panic("implement me")
}

func New(db postgresql.Client) planets.IReadRepository {
	r := new(Repositoty)
	r.client = db
	return r
}

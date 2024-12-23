package repository

import (
	"context"
	"example/internal/features/peoples"

	"github.com/ihatiko/chef/components/clients/redis"
)

type repository struct {
	client redis.Client
}

func (r repository) Cache(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func New(client redis.Client) peoples.IRepository {
	return &repository{
		client: client,
	}
}

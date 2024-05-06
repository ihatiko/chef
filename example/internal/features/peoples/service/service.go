package service

import (
	"context"
	"example/internal/features/peoples"
)

type service struct {
	repository peoples.IRepository
}

func (s service) Cache(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func New(repository peoples.IRepository) peoples.IService {
	return &service{repository: repository}
}

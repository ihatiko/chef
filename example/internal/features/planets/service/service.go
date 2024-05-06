package service

import (
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type service struct {
	readRepository planets.IReadRepository
}

func New(readRepository planets.IReadRepository) planets.IService {
	return &service{readRepository: readRepository}
}

func (s service) Load(request daemon.Request) error {
	//TODO implement me
	panic("implement me")
}

func (s service) Update(request cron.Request) error {
	//TODO implement me
	panic("implement me")
}

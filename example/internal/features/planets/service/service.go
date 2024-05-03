package service

import (
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type Service struct {
	readRepository planets.IReadRepository
}

func New(readRepository planets.IReadRepository) planets.IService {
	return &Service{readRepository: readRepository}
}

func (s Service) Load(request daemon.Request) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) Update(request cron.Request) error {
	//TODO implement me
	panic("implement me")
}

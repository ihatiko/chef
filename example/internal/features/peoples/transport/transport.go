package transport

import (
	"context"
	"example/internal/features/peoples"
	protoPeople "example/pkg/protoc/peoples"
)

type transport struct {
	service peoples.IService
}

func (t transport) UpdatePeople(ctx context.Context, request *protoPeople.UpdatePeopleRequest) (*protoPeople.UpdatePeopleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func New(service peoples.IService) peoples.ITransport {
	return &transport{service: service}
}

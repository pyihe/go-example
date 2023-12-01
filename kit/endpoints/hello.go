package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pyihe/go-example/kit/service"
)

type HTTPEndpoints struct {
	httpSrv service.HelloService
	grpcSrv service.VisitService
}

func MakeEndpoints() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

	}
}

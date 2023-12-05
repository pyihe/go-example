package endpoints

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pyihe/go-example/kit/internal/middleware"
	"github.com/pyihe/go-example/kit/protocol"
	"github.com/pyihe/go-example/kit/service"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

type EndpointSet struct {
	HelloEndpoint endpoint.Endpoint
	VisitEndpoint endpoint.Endpoint
}

func NewEndpointSet(srv *service.Service, logger log.Logger, tracer stdopentracing.Tracer) EndpointSet {
	var helloEndpoint = wrapEndpoint("Hello", makeHelloEndpoint(srv), logger, tracer)
	var visitEndpoint = wrapEndpoint("Visit", makeVisitEndpoint(srv), logger, tracer)

	return EndpointSet{
		HelloEndpoint: helloEndpoint,
		VisitEndpoint: visitEndpoint,
	}
}

func wrapEndpoint(operationName string, ep endpoint.Endpoint, logger log.Logger, tracer stdopentracing.Tracer) endpoint.Endpoint {
	ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(ep)
	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
	ep = opentracing.TraceServer(tracer, operationName)(ep)
	ep = middleware.LoggingMiddleware(logger)(ep)
	return ep
}

func makeHelloEndpoint(s *service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*protocol.HelloRequest)
		if req == nil || !ok {
			return nil, ErrInvalidRequest
		}
		resp, err := s.Hello(req.GetName())
		return &protocol.HelloResponse{Echo: resp}, err
	}
}

func makeVisitEndpoint(s *service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*protocol.VisitRequest)
		if req == nil || !ok {
			return nil, ErrInvalidRequest
		}
		resp, err := s.Visit(req.Place)
		return &protocol.VisitResponse{Echo: resp}, err
	}
}

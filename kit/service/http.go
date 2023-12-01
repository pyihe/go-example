package service

import (
	"context"
	"fmt"
)

type HelloService interface {
	Greet(ctx context.Context, name string) (string, error)
}

func NewHttpService() HelloService {
	srv := &httpService{}

	return srv
}

type httpService struct{}

func (s *httpService) Greet(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("hello %s, this is your http service.", name), nil
}

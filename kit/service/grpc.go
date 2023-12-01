package service

import "fmt"

type VisitService interface {
	Visit(place string) (string, error)
}

func NewGRPCService() VisitService {
	srv := &gRPCService{}

	return srv
}

type gRPCService struct{}

func (s *gRPCService) Visit(place string) (string, error) {
	return fmt.Sprintf("hello, you are visiting %s.", place), nil
}

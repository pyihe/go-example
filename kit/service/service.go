package service

import (
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Hello(name string) (string, error) {
	return fmt.Sprintf("hello %s.", name), nil
}

func (s *Service) Visit(place string) (string, error) {
	return fmt.Sprintf("wellcome to %s.", place), nil
}

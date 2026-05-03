package users

import "fmt"

type Service interface {
	Create(request CreateUserRequest) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s service) Create(request CreateUserRequest) error {
	fmt.Println(request)
	fmt.Println("Create user service")
	return nil
}

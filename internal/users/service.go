package users

import (
	"log"
)

type Service interface {
	Create(request CreateUserRequest) error
}

type service struct {
	log  *log.Logger
	repo Reposiotry
}

func NewService(log *log.Logger, repo Reposiotry) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(request CreateUserRequest) error {
	s.log.Println(request)
	s.log.Println("Create user service")
	s.repo.Create(&User{})
	return nil
}

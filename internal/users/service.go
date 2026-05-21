package users

import (
	"courses/internal/domain"
	"fmt"
	"log"
)

type (
	Filters struct {
		Firtsname string
		Email     string
	}

	Service interface {
		Create(request CreateUserRequest) (*domain.User, error)
		GetAll(filter Filters, offset, limit int64) ([]domain.User, error)
		GetById(uuid string) (*domain.User, error)
		ExistsById(uuid string) error
		Delete(uuid string) error
		Update(uuid string, request CreateUserRequest) (*domain.User, error)
		Count(filter Filters) (int64, error)
	}
	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(request CreateUserRequest) (*domain.User, error) {
	if s.repo.ExistsByEmail(request.Email) {
		return nil, fmt.Errorf("Already exist a user with that email %s", request.Email)
	}
	user := domain.User{
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Email:     request.Email,
		Phone:     request.Phone,
	}
	userCreated, err := s.repo.Create(&user)
	if err != nil {
		return nil, err
	}
	return userCreated, nil
}

func (s service) GetAll(filter Filters, offset, limit int64) ([]domain.User, error) {
	usersFound, err := s.repo.GetAll(filter, offset, limit)
	if err != nil {
		return nil, err
	}
	if len(usersFound) == 0 {
		return nil, fmt.Errorf("Users not found")
	}
	return usersFound, nil
}

func (s service) ExistsById(uuid string) error {
	existsUser := s.repo.ExistsById(uuid)
	if !existsUser {
		return fmt.Errorf("Not exists an user with that id %s", uuid)
	}
	return nil
}

func (s service) GetById(uuid string) (*domain.User, error) {
	userFound, err := s.repo.GetById(uuid)
	if err != nil {
		return nil, err
	}
	if userFound == nil {
		return nil, fmt.Errorf("Not exists an user with that id %s", uuid)
	}
	return userFound, nil
}

func (s service) Delete(uuid string) error {
	existsUser := s.repo.ExistsById(uuid)
	if !existsUser {
		return fmt.Errorf("Not exists a user with that id %s", uuid)
	}
	result := s.repo.Delete(uuid)
	if result != nil {
		return result
	}
	return nil
}

func (s *service) Update(uuid string, request CreateUserRequest) (*domain.User, error) {
	existsUser := s.repo.ExistsById(uuid)
	if !existsUser {
		return nil, fmt.Errorf("Not exists a user with that id %s", uuid)
	}
	user := domain.User{
		Id:        uuid,
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Email:     request.Email,
		Phone:     request.Phone,
	}
	userUpdated, err := s.repo.Update(&user)
	if err != nil {
		return nil, err
	}
	return userUpdated, nil
}

func (s *service) Count(filter Filters) (int64, error) {
	return s.repo.Count(filter)
}

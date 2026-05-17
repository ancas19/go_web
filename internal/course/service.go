package course

import (
	"fmt"
	"log"
	"time"
)

type (
	Service interface {
		Create(courseRequest CreateCourseReq) (*Course, error)
		Delete(idCourse string) error
		GetById(idCourse string) (*Course, error)
		Count(name string) (int64, error)
		GetAllCourses(name string, page, limit int64) ([]Course, error)
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

func (s service) Create(courseRequest CreateCourseReq) (*Course, error) {
	startDateparsed, err := time.Parse("2006-01-02", courseRequest.StartDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	endDateParsed, err := time.Parse("2006-01-02", courseRequest.EndtDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	courseToCreate := Course{
		Name:      courseRequest.Name,
		EndDate:   endDateParsed,
		StartDate: startDateparsed,
	}
	courseCreated, err := s.repo.Create(&courseToCreate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return courseCreated, nil
}

func (s service) Delete(idCourse string) error {
	existsCourse, err := s.repo.ExistsById(idCourse)
	if err != nil {
		return err
	}
	if !existsCourse {
		return fmt.Errorf("Not exists a course with that id %s", idCourse)
	}
	err = s.repo.Delete(idCourse)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetById(idCourse string) (*Course, error) {
	existsCourse, err := s.repo.ExistsById(idCourse)
	if err != nil {
		return nil, err
	}
	if !existsCourse {
		return nil, fmt.Errorf("No exists a course with that id %s", idCourse)
	}
	courseFound, err := s.repo.GetById(idCourse)
	if err != err {
		return nil, err
	}
	return courseFound, nil
}

func (s *service) GetAllCourses(name string, offset, limit int64) ([]Course, error) {
	usersFound, err := s.repo.GetAllCourses(name, offset, limit)
	if err != nil {
		return nil, err
	}
	if len(usersFound) == 0 {
		return nil, fmt.Errorf("Users not found")
	}
	return usersFound, nil
}

func (s *service) Count(name string) (int64, error) {
	count, err := s.repo.Count(name)
	if err != nil {
		return 0, err
	}
	return count, nil
}

package course

import (
	"courses/internal/domain"
	"fmt"
	"log"
	"time"
)

type (
	Service interface {
		Create(courseRequest CreateCourseReq) (*domain.Course, error)
		Delete(idCourse string) error
		GetById(idCourse string) (*domain.Course, error)
		Count(name string) (int64, error)
		GetAllCourses(name string, page, limit int64) ([]domain.Course, error)
		Update(idCourse string, courseRequest CreateCourseReq) (*domain.Course, error)
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

func (s service) Create(courseRequest CreateCourseReq) (*domain.Course, error) {
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
	courseToCreate := domain.Course{
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

func (s *service) GetById(idCourse string) (*domain.Course, error) {
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

func (s *service) GetAllCourses(name string, offset, limit int64) ([]domain.Course, error) {
	usersFound, err := s.repo.GetAllCourses(name, offset, limit)
	if err != nil {
		return nil, err
	}
	if len(usersFound) == 0 {
		return nil, fmt.Errorf("Users not found")
	}
	return usersFound, nil
}

func (s *service) Update(idCourse string, courseRequest CreateCourseReq) (*domain.Course, error) {
	if idCourse == "" {
		return nil, fmt.Errorf("course id cannot be empty")
	}
	existsCourse, err := s.repo.ExistsById(idCourse)
	if err != nil {
		return nil, fmt.Errorf("error fetching course: %w", err)
	}
	if !existsCourse {
		return nil, fmt.Errorf("course with id %s not found", idCourse)
	}
	startDate, err := s.mapDate(courseRequest.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := s.mapDate(courseRequest.EndtDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}
	infoToUpdate := findInfoToUpdate(courseRequest.Name, startDate, endDate)
	if len(infoToUpdate) == 0 {
		return nil, fmt.Errorf("There aren't iformation to update")
	}
	courseUpdated, err := s.repo.Update(idCourse, infoToUpdate)
	if err != nil {
		return nil, fmt.Errorf("error updating course: %w", err)
	}

	return courseUpdated, nil
}

func (s *service) Count(name string) (int64, error) {
	count, err := s.repo.Count(name)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *service) mapDate(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}
	dateParsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return &dateParsed, nil
}

func findInfoToUpdate(name string, startDate, endDate *time.Time) map[string]any {
	infoToUpdate := make(map[string]any)
	if name != "" {
		infoToUpdate["name"] = name
	}
	if startDate != nil {
		infoToUpdate["start_date"] = *startDate
	}
	if startDate != nil {
		infoToUpdate["end_date"] = *startDate
	}
	return infoToUpdate
}

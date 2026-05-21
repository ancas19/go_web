package enrollment

import (
	"courses/internal/course"
	"courses/internal/domain"
	"courses/internal/users"
	"log"
)

type (
	Service interface {
		Create(enrollment EnrollmentRequest) (*domain.Enrollment, error)
	}

	service struct {
		log           *log.Logger
		userService   users.Service
		courseService course.Service
		repo          Repository
	}
)

func NewService(log *log.Logger, repo Repository, userService users.Service, courseService course.Service) Service {
	return &service{
		log:           log,
		repo:          repo,
		userService:   userService,
		courseService: courseService,
	}
}

func (s *service) Create(enrollment EnrollmentRequest) (*domain.Enrollment, error) {
	err := s.courseService.ExistsById(enrollment.CourseId)
	if err != nil {
		return nil, err
	}
	err = s.userService.ExistsById(enrollment.UserId)
	if err != nil {
		return nil, err
	}
	enrollmentToCreate := domain.Enrollment{
		UserId:   enrollment.UserId,
		CourseId: enrollment.CourseId,
		Status:   "P",
	}
	enrollmentCreated, err := s.repo.Create(&enrollmentToCreate)
	if err != nil {
		return nil, err
	}
	return enrollmentCreated, nil
}

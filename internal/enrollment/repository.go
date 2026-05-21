package enrollment

import (
	"courses/internal/domain"
	"log"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enrollment *domain.Enrollment) (*domain.Enrollment, error)
	}
	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(enrollment *domain.Enrollment) (*domain.Enrollment, error) {
	if err := repo.db.Create(enrollment).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}
	return enrollment, nil
}

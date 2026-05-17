package course

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) (*Course, error)
		Delete(idCourse string) error
		ExistsById(idCourse string) (bool, error)
		GetById(idCourse string) (*Course, error)
		Count(name string) (int64, error)
		GetAllCourses(name string, page, limit int64) ([]Course, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (repo *repo) Create(course *Course) (*Course, error) {
	if err := repo.db.Create(course).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}
	return course, nil
}

func (repo *repo) ExistsById(idCourse string) (bool, error) {
	var count int64
	result := repo.db.Model(&Course{}).Where("id = ?", idCourse).Count(&count)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return false, result.Error
	}
	return count > 0, nil
}

func (repo *repo) Delete(idCourse string) error {
	course := Course{Id: idCourse}
	result := repo.db.Delete(&course)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}
	return nil
}

func (repo *repo) GetById(idCourse string) (*Course, error) {
	var courseFound Course = Course{Id: idCourse}
	result := repo.db.First(&courseFound)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	return &courseFound, nil
}

func (repo *repo) Count(name string) (int64, error) {
	var count int64
	tx := repo.db.Model(Course{})
	tx = applyFilters(tx, name)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *repo) GetAllCourses(name string, offset, limit int64) ([]Course, error) {
	var courses []Course
	ctx := repo.db.Model(Course{})
	ctx = applyFilters(ctx, name)
	ctx.Limit(int(limit))
	ctx.Offset(int(offset))
	result := ctx.Order("created_at desc").Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}

func applyFilters(tx *gorm.DB, name string) *gorm.DB {
	if name != "" {
		name = fmt.Sprintf("%%%s%%", strings.ToLower(name))
		tx = tx.Where("lower(first_name) like ?", name)
	}
	return tx
}

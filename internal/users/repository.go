package users

import (
	"courses/internal/domain"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) (*domain.User, error)
	ExistsByEmail(email string) bool
	GetAll(filter Filters, offset, limit int64) ([]domain.User, error)
	GetById(uuid string) (*domain.User, error)
	Delete(uuid string) error
	ExistsById(uuid string) bool
	Update(User *domain.User) (*domain.User, error)
	Count(filter Filters) (int64, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepos(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *domain.User) (*domain.User, error) {
	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}
	return user, nil
}

func (repo *repo) ExistsByEmail(email string) bool {
	var count int64
	repo.db.Model(&domain.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (repo *repo) GetAll(filter Filters, offset, limit int64) ([]domain.User, error) {
	var users []domain.User
	ctx := repo.db.Model(&domain.User{})
	ctx = applyFilters(ctx, filter)
	ctx.Limit(int(limit))
	ctx.Offset(int(offset))
	result := ctx.Order("created_at desc").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *repo) GetById(uuid string) (*domain.User, error) {
	var userFound domain.User = domain.User{Id: uuid}
	result := repo.db.First(&userFound)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &userFound, nil
}

func (repo *repo) Delete(uuid string) error {
	user := domain.User{Id: uuid}
	result := repo.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *repo) ExistsById(uuid string) bool {
	var count int64
	repo.db.Model(&domain.User{}).Where("id = ?", uuid).Count(&count)
	return count > 0
}

func (repo *repo) Update(user *domain.User) (*domain.User, error) {
	infoToUpdate := make(map[string]interface{})
	if user.Firstname != "" {
		infoToUpdate["first_name"] = user.Firstname
	}
	if user.Lastname != "" {
		infoToUpdate["last_name"] = user.Lastname
	}
	if user.Email != "" {
		infoToUpdate["email"] = user.Email
	}
	if user.Phone != "" {
		infoToUpdate["phone"] = user.Phone
	}
	result := repo.db.Model(&domain.User{}).Where("id = ?", user.Id).Updates(infoToUpdate)
	if result.Error != nil {
		return nil, result.Error
	}
	var updatedUser domain.User
	if err := repo.db.Where("id = ?", user.Id).First(&updatedUser).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *repo) Count(filters Filters) (int64, error) {
	var count int64
	tx := r.db.Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Firtsname != "" {
		filters.Firtsname = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Firtsname))
		tx = tx.Where("lower(first_name) like ?", filters.Firtsname)
	}
	if filters.Email != "" {
		filters.Firtsname = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Firtsname))
		tx = tx.Where("lower(email) like ?", filters.Email)
	}
	return tx
}

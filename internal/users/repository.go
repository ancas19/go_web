package users

import (
	"log"

	"gorm.io/gorm"
)

type Reposiotry interface {
	Create(user *User) error
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepos(log *log.Logger, db *gorm.DB) Reposiotry {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *User) error {
	repo.log.Println(user)
	return nil
}

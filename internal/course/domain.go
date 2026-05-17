package course

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	Id        string         `json:"id" gorm:"type:varchar(36);not null;primary_key;unique_index"`
	Name      string         `json:"name" gorm:"type:varchar(60);not null"`
	StartDate time.Time      `json:"startDate"`
	EndDate   time.Time      `json:"endDate"`
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	Deleted   gorm.DeletedAt `json:"-"`
}

func (u *Course) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}

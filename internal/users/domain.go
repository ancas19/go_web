package users

import "time"

type User struct {
	Id        string     `json:"id" gorm:"type:varchar(36);not null;primary_key;unique_index"`
	Firstname string     `json:"firstname" gorm:"type:varchar(60);not null"`
	Lastname  string     `json:"lastname" gorm:"type:varchar(60);not null"`
	Email     string     `json:"email" gorm:"type:varchar(50);not null"`
	Phone     string     `json:"phone"  gorm:"type:varchar(20);not null"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

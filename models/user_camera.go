package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Pseudo   string `gorm:"unique;not null" json:"pseudo"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

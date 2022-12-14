package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Pseudo   string `gorm:"unique;not null" json:"pseudo"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

/*
func (u *User) RegisterCamera(db *gorm.DB, c *Camera) error {
	return db.Model(&u).Association("Cameras").Append(c)
}

func (u *User) DeleteCamera(db *gorm.DB, c *Camera) error {
	return db.Model(&u).Association("Cameras").Delete(c)
}
*/

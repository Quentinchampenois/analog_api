package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Pseudo   string   `gorm:"unique;not null" json:"pseudo"`
	Password string   `json:"password"`
	Role     string   `json:"role"`
	Cameras  []Camera `gorm:"many2many:user_cameras;"`
}

func (u *User) registerCamera(db *gorm.DB, c *Camera) error {
	return db.Model(&u).Association("Cameras").Append(c)
}

package models

import (
	"gorm.io/gorm"
	"time"
)

type UserCameraFilm struct {
	gorm.Model
	UserID    int       `json:"user_id"`
	User      User      `json:"-"`
	CameraID  int       `json:"-"`
	Camera    Camera    `json:"camera"`
	FilmID    int       `json:"-"`
	Film      Film      `json:"film"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

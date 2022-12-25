package models

import (
	"gorm.io/gorm"
	"time"
)

type UserCameraFilm struct {
	gorm.Model
	UserCameraID int       `json:"user_camera_id"`
	FilmID       int       `json:"film_id""`
	StartDate    time.Time `gorm:"not null" json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

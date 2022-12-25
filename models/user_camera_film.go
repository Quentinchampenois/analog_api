package models

import (
	"gorm.io/gorm"
)

type UserCameraFilm struct {
	gorm.Model
	UserCameraID int   `json:"user_camera_id"`
	FilmID       int   `json:"film_id""`
	StartDate    int64 `gorm:"not null" json:"start_date"`
	EndDate      int64 `json:"end_date"`
}

package models

import (
	"gorm.io/gorm"
)

type UserCameraFilm struct {
	gorm.Model
	UserCameraID int        `json:"user_camera_id"`
	UserCamera   UserCamera `json:"-"`
	FilmID       int        `json:"-"`
	Film         Film       `json:"film"`
	StartDate    int64      `gorm:"not null" json:"start_date"`
	EndDate      int64      `json:"end_date"`
}

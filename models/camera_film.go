package models

import (
	"gorm.io/gorm"
	"time"
)

type CameraFilm struct {
	gorm.Model
	UserID    int       `json:"-"`
	User      User      `json:"user"`
	CameraID  int       `json:"-"`
	Camera    Camera    `json:"camera"`
	FilmID    int       `json:"-"`
	Film      Film      `json:"film_id"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

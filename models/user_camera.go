package models

import (
	"gorm.io/gorm"
)

type UserCamera struct {
	gorm.Model

	ID       int `json:"ID"`
	UserID   int `json:"user_id"`
	CameraID int `json:"camera_id"`
}

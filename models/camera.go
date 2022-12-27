package models

import (
	"gorm.io/gorm"
)

type Camera struct {
	gorm.Model

	ID     int    `json:"ID"`
	Name   string `json:"name"`
	TypeID int    `json:"type_id"`
	Type   Type   `json:"type"`
	Focus  string `json:"focus"`
}

func GetCameras(db *gorm.DB, start, count int) ([]Camera, error) {
	var cameras []Camera

	if err := db.Limit(count).Offset(start).Preload("Type").Find(&cameras).Error; err != nil {
		return nil, err
	}

	return cameras, nil
}

func (c *Camera) CreateCamera(db *gorm.DB) bool {
	var t Type
	if !t.GetType(db, c.TypeID) {
		return false
	}

	camera := &Camera{
		ID:     c.ID,
		Name:   c.Name,
		Focus:  c.Focus,
		TypeID: t.ID,
		Type: Type{
			ID:   t.ID,
			Name: t.Name,
		},
	}
	return db.Create(&camera).Error == nil
}

func (c *Camera) GetCamera(db *gorm.DB, id int) bool {
	db.Preload("Type").Find(&c, id)

	return !(c.ID == 0)
}

func (c *Camera) UpdateCamera(db *gorm.DB) error {
	return db.Save(&c).Error
}

func (c *Camera) DeleteCamera(db *gorm.DB) error {
	return db.Delete(&c).Error
}

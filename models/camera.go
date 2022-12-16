package models

import (
	"gorm.io/gorm"
	"log"
)

type Camera struct {
	gorm.Model

	ID     int    `json:"ID"`
	Name   string `json:"name"`
	TypeID int    `json:"-"`
	Type   Type   `json:"type"`
	Focus  string `json:"focus"`
}

func GetCameras(db *gorm.DB, start, count int) ([]Camera, error) {
	var cameras []Camera

	if err := db.Preload("Type").Find(&cameras).Error; err != nil {
		log.Fatalf("Error append in getCameras : \n%v\n", err)
		return nil, nil
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
	db.Create(camera)
	db.Save(camera)
	return true
}

func (c *Camera) GetCamera(db *gorm.DB, id int) bool {
	db.Preload("Type").Find(&c, id)

	if c.ID == 0 {
		return false
	}
	return true
}

func (c *Camera) UpdateCamera(db *gorm.DB) error {
	err := db.Save(&c).Error
	return err
}

func (c *Camera) DeleteCamera(db *gorm.DB) error {
	err := db.Delete(&c).Error
	return err
}

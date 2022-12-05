package main

import (
	"gorm.io/gorm"
	"log"
)

// camera - Define a camera with specificity
// TODO - Change Type for relation
// TODO - Change Focus for relation
type Camera struct {
	gorm.Model

	ID     int    `json:"ID"`
	Name   string `json:"name"`
	TypeID int    `json:"typeID"`
	Type   Type   `json:"type"`
	Focus  string `json:"focus"`
	Film   int    `json:"film"`
}

type Type struct {
	gorm.Model

	ID   int    `json:"ID"`
	Name string `json:"name"`
}

func getCameras(db *gorm.DB, start, count int) ([]Camera, error) {
	var cameras []Camera

	if err := db.Preload("Type").Find(&cameras).Error; err != nil {
		log.Fatalf("Error append in getCameras : \n%v\n", err)
		return nil, nil
	}

	return cameras, nil
}

func (c *Camera) createCamera(db *gorm.DB) bool {
	var t Type
	if !t.getType(db, c.TypeID) {
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
		Film: c.Film,
	}
	db.Create(camera)
	db.Save(camera)
	return true
}

func (c *Camera) getCamera(db *gorm.DB, id int) bool {
	db.Preload("Type").Find(&c, id)

	if c.ID == 0 {
		return false
	}
	return true
}

func (c *Camera) updateCamera(db *gorm.DB) error {
	err := db.Save(&c).Error
	return err
}

func (c *Camera) deleteCamera(db *gorm.DB) error {
	err := db.Delete(&c).Error
	return err
}

func getTypes(db *gorm.DB, start, count int) ([]Type, error) {
	var cameras []Type

	if err := db.Find(&cameras).Error; err != nil {
		log.Fatalf("Error append in getTypes : \n%v\n", err)
		return nil, nil
	}

	return cameras, nil
}

func (t *Type) createType(db *gorm.DB) {
	db.Create(&t)
}

func (t *Type) getType(db *gorm.DB, id int) bool {
	db.First(&t, id)

	if t.ID == 0 {
		return false
	}
	return true
}

func (t *Type) updateType(db *gorm.DB) {
	db.Save(&t)
}

func (t *Type) deleteType(db *gorm.DB) {
	db.Delete(&t)
}

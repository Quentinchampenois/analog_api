package main

import (
	"gorm.io/gorm"
	"log"
)

// camera - Define a camera with specificity
// TODO - Change Type for relation
// TODO - Change Focus for relation
type camera struct {
	gorm.Model

	ID    int    `json:"ID"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Focus string `json:"focus"`
	Film  int    `json:"film"`
}

func getCameras(db *gorm.DB, start, count int) ([]camera, error) {
	var cameras []camera

	if err := db.Find(&cameras).Error; err != nil {
		log.Fatalf("Error append in getCameras : \n%v\n", err)
		return nil, nil
	}

	return cameras, nil
}

func (c *camera) createCamera(db *gorm.DB) {
	db.Create(&c)
}

func (c *camera) getCamera(db *gorm.DB, id int) bool {
	db.Find(&c, id)

	if c.ID == 0 {
		return false
	}
	return true
}

func (c *camera) updateCamera(db *gorm.DB) {
	db.Save(&c)
}

func (c *camera) deleteCamera(db *gorm.DB) {
	db.Delete(&c)
}

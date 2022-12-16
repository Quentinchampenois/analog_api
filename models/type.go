package models

import (
	"gorm.io/gorm"
	"log"
)

type Type struct {
	gorm.Model

	ID   int    `json:"ID"`
	Name string `json:"name"`
}

func GetTypes(db *gorm.DB, start, count int) ([]Type, error) {
	var types []Type

	if err := db.Find(&types).Error; err != nil {
		log.Fatalf("Error append in getTypes : \n%v\n", err)
		return nil, nil
	}

	return types, nil
}

func (t *Type) CreateType(db *gorm.DB) {
	db.Create(&t)
}

func (t *Type) GetType(db *gorm.DB, id int) bool {
	db.First(&t, id)

	if t.ID == 0 {
		return false
	}
	return true
}

func (t *Type) UpdateType(db *gorm.DB) {
	db.Save(&t)
}

func (t *Type) DeleteType(db *gorm.DB) {
	db.Delete(&t)
}

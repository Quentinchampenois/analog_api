package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

type Film struct {
	gorm.Model

	ID    int    `json:"ID"`
	Name  string `json:"name"`
	Size  int    `json:"size"`
	Color bool   `json:"color"`
	Brand string `json:"brand"`
}

func GetFilms(db *gorm.DB, start, count int) ([]Film, error) {
	var films []Film

	if err := db.Find(&films).Error; err != nil {
		log.Fatalf("Error append in getFilms : \n%v\n", err.Error())
		return nil, nil
	}

	return films, nil
}

func (f *Film) CreateFilm(db *gorm.DB) bool {
	film := &Film{
		ID:    f.ID,
		Name:  f.Name,
		Size:  f.Size,
		Color: f.Color,
		Brand: f.Brand,
	}

	if err := db.Create(film).Error; err != nil {
		fmt.Println(err.Error())
		return false
	}

	if err := db.Save(film).Error; err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func (f *Film) GetFilm(db *gorm.DB, id int) error {
	if err := db.First(&f, id).Error; err != nil {
		return err
	}

	return nil
}

func (f *Film) UpdateFilm(db *gorm.DB) {
	db.Save(&f)
}

func (f *Film) DeleteFilm(db *gorm.DB) error {
	if err := db.Delete(&f).Error; err != nil {
		return err
	}

	return nil
}

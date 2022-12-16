package main

import (
	"github.com/quentinchampenois/analog_api/models"
	"log"
)

func (a *App) migrate() {
	if err := a.DB.AutoMigrate(&models.Camera{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.Type{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.Film{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.CameraFilm{}); err != nil {
		log.Fatal(err)
	}
}

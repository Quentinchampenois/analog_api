package main

import (
	"fmt"
	"github.com/quentinchampenois/analog_api/models"
	"log"
	"time"
)

func (a *App) migrate() {
	if err := a.DB.AutoMigrate(&models.Type{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.Camera{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.Film{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	if err := a.DB.AutoMigrate(&models.UserCameraFilm{}); err != nil {
		log.Fatal(err)
	}

	var cameras []models.Camera
	a.DB.Find(&cameras)

	if len(cameras) <= 0 {
		a.seed()
	}
}

func (a *App) seed() {
	fmt.Println("Seeding Database")

	types := []models.Type{
		{
			Name: "Compact Point & Shoot",
		},
		{
			Name: "Télémétrique",
		},
		{
			Name: "Reflex",
		},
	}

	a.DB.Create(&types)

	cameras := []models.Camera{{
		Name:   "Fujica DL-100",
		TypeID: 1,
		Focus:  "Autofocus",
	}, {
		Name:   "Minolta AF-S",
		TypeID: 1,
		Focus:  "Autofocus",
	}, {
		Name:   "Minolta AF-S 2",
		TypeID: 1,
		Focus:  "Autofocus",
	}, {
		Name:   "Télémétrique",
		TypeID: 2,
		Focus:  "Autofocus",
	}, {
		Name:   "Reflex",
		TypeID: 2,
		Focus:  "Autofocus",
	},
	}

	a.DB.Create(&cameras)

	encryptedPwd, err := a.EncryptedPassword("password")
	if err != nil {
		log.Fatalln(err.Error())
	}
	users := []models.User{{
		Pseudo:   "user",
		Password: encryptedPwd,
		Role:     "user",
	},
		{
			Pseudo:   "admin",
			Password: encryptedPwd,
			Role:     "admin",
		}}

	a.DB.Create(users)

	userCameras := []models.UserCamera{{
		UserID:   1,
		CameraID: 2,
	},
		{
			UserID:   1,
			CameraID: 1,
		},
		{
			UserID:   1,
			CameraID: 3,
		},
		{
			UserID:   3,
			CameraID: 3,
		}}

	a.DB.Create(&userCameras)

	films := []models.Film{{
		Name:  "APX",
		Size:  135,
		Iso:   400,
		Color: false,
		Brand: "APX",
	}, {
		Name:  "Color Plus",
		Size:  135,
		Iso:   200,
		Color: true,
		Brand: "Kodak",
	}, {
		Name:  "Gold",
		Size:  135,
		Iso:   200,
		Color: true,
		Brand: "Kodak",
	}}

	a.DB.Create(&films)

	userCameraFilms := []models.UserCameraFilm{{
		UserCameraID: 1,
		FilmID:       1,
		StartDate:    time.Date(2021, time.October, 10, 12, 0, 0, 0, time.UTC).Unix(),
		EndDate:      time.Date(2021, time.November, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}, {
		UserCameraID: 1,
		FilmID:       2,
		StartDate:    time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
		EndDate:      time.Date(2021, time.November, 20, 23, 0, 0, 0, time.UTC).Unix(),
	}, {
		UserCameraID: 1,
		FilmID:       3,
		StartDate:    time.Now().Unix(),
		EndDate:      0,
	}}

	a.DB.Create(&userCameraFilms)
}

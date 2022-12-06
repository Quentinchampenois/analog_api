package main

import (
	"github.com/quentinchampenois/analog_api/models"
	"log"
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	if err := a.DB.AutoMigrate(&Camera{}); err != nil {
		log.Fatal(err)
	}
	if err := a.DB.AutoMigrate(&Type{}); err != nil {
		log.Fatal(err)
	}
	if err := a.DB.AutoMigrate(&models.Film{}); err != nil {
		log.Fatal(err)
	}

	a.Run(":8080")
}

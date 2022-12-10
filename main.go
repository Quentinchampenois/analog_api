package main

import (
	"github.com/joho/godotenv"
	"github.com/quentinchampenois/analog_api/models"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	a := App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

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

	a.Run(":8080")
}

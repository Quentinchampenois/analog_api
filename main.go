package main

import (
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

	if err := a.DB.AutoMigrate(&camera{}); err != nil {
		log.Fatal(err)
	}

	a.Run(":8080")
}

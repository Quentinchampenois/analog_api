package main

import (
	"github.com/joho/godotenv"
	"github.com/quentinchampenois/analog_api/configs"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	a := App{
		Configs: configs.Config{
			Server:   configs.Server{},
			Database: configs.Database{},
		},
	}
	a.Initialize()
	a.migrate()
	a.Run()
}

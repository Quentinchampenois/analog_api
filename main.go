package main

import (
	"fmt"
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

	const tableCreationQuery = `CREATE TABLE IF NOT EXISTS cameras
(
    id SERIAL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    focus TEXT NOT NULL,
    film INTEGER NOT NULL,
    CONSTRAINT cameras_pkey PRIMARY KEY (id)
)`

	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initialized application")
	a.Run(":8080")
}

package configs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type Database struct {
	username string
	password string
	dbName   string
	host     string
	port     string
}

func (d *Database) Setup() {
	var missingEnvVar []string

	d.host = d.getHost()
	d.port = d.getPort()
	if d.username = d.getEnvOrFail("DB_USERNAME"); d.username == "" {
		missingEnvVar = append(missingEnvVar, "DB_USERNAME")
	}
	if d.password = d.getEnvOrFail("DB_PASSWORD"); d.password == "" {
		missingEnvVar = append(missingEnvVar, "DB_PASSWORD")
	}
	if d.dbName = d.getEnvOrFail("DB_NAME"); d.dbName == "" {
		missingEnvVar = append(missingEnvVar, "DB_NAME")
	}

	if len(missingEnvVar) > 0 {
		for _, val := range missingEnvVar {
			fmt.Printf("Missing '%v'\n", val)
		}

		log.Fatalln("Please add missing required environment variables")
	}
}

func (d *Database) Connection() *gorm.DB {
	dsn :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/paris", d.host, d.port, d.username, d.password, d.dbName)

	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return gormDb
}

func (d *Database) getHost() string {
	if host := os.Getenv("DB_HOST_NAME"); host != "" {
		return host
	}

	return "localhost"
}

func (d *Database) getPort() string {
	if port := os.Getenv("DB_PORT"); port != "" {
		return port
	}

	return "5432"
}

func (d *Database) getEnvOrFail(key string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}

	return ""
}

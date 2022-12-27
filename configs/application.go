package configs

import (
	"log"
	"os"
	"strconv"
)

type Application struct {
	PaginateLimit int
}

func (a *Application) Setup() {
	a.PaginateLimit = a.getPagination()
}

func (a *Application) getPagination() int {
	if pagination := os.Getenv("APP_PAGINATION"); pagination != "" {
		pg, err := strconv.Atoi(pagination)
		if err != nil {
			log.Println("Argument 'APP_PAGINATION' is not an integer, applying default value : 20")
			return 20
		}

		return pg
	}

	return 20
}

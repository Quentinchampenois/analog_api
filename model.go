package main

import (
	"database/sql"
	"log"
)

// camera - Define a camera with specificity
// TODO - Change Type for relation
// TODO - Change Focus for relation
type camera struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Focus string `json:"focus"`
	Film  int    `json:"film"`
}

func (c *camera) getCamera(db *sql.DB) error {
	err := db.QueryRow("SELECT id, name, type, focus, film FROM cameras WHERE id=$1", c.ID).Scan(&c.ID, &c.Name, &c.Type, &c.Focus, &c.Film)
	if err != nil {
		return err
	}
	return nil
}

func (c *camera) updateCamera(db *sql.DB) error {
	_, err := db.Exec("UPDATE cameras SET name=$1, type=$2, focus=$3, film=$4 WHERE ID=$5", &c.Name, &c.Type, &c.Focus, &c.Film, &c.ID)

	return err
}

func (c *camera) deleteCamera(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM cameras WHERE ID=$1", &c.ID)

	return err
}

func (c *camera) createCamera(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO cameras(name, type, focus, film) VALUES($1, $2, $3, $4) RETURNING id",
		c.Name, c.Type, c.Focus, c.Film).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func getCameras(db *sql.DB, start, count int) ([]camera, error) {
	rows, err := db.Query("SELECT id, name, type, focus, film FROM cameras LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error append in getCameras : \n%v\n", err)
			return
		}
	}(rows)

	var cameras []camera
	for rows.Next() {
		var c camera
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Focus, &c.Film); err != nil {
			return nil, err
		}
		cameras = append(cameras, c)
	}
	return cameras, nil
}

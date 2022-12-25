package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/quentinchampenois/analog_api/models"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (a *App) getCameras(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 20 || count < 1 {
		count = a.Configs.Application.PaginateLimit
	}

	if start < 0 {
		start = 0
	}

	cameras, err := models.GetCameras(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cameras)
}

func (a *App) createCamera(w http.ResponseWriter, r *http.Request) {
	var c models.Camera
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println(err)
		}
	}(r.Body)

	if !c.CreateCamera(a.DB) {
		respondWithError(w, http.StatusNotFound, "ID of type not found")
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) getCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	var c models.Camera
	if !c.GetCamera(a.DB, id) {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) updateCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var c models.Camera
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	c.ID = id
	if err := c.UpdateCamera(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Camera ID")
		return
	}

	c := models.Camera{ID: id}
	if err := c.DeleteCamera(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

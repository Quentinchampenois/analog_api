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

func (a *App) getFilms(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	films, err := models.GetFilms(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, films)
}

func (a *App) createFilm(w http.ResponseWriter, r *http.Request) {
	var f models.Film
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&f); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}(r.Body)

	if !f.CreateFilm(a.DB) {
		respondWithError(w, http.StatusNotFound, "ID of type not found")
		return
	}

	respondWithJSON(w, http.StatusCreated, f)
}

func (a *App) getFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	var f models.Film
	if err = f.GetFilm(a.DB, id); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, f)
}

func (a *App) updateFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var f models.Film

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&f); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	f.ID = id
	var update models.Film
	if err = update.GetFilm(a.DB, f.ID); err != nil {
		respondWithError(w, http.StatusNotFound, "Film does not exist")
		return
	}

	f.UpdateFilm(a.DB)
	respondWithJSON(w, http.StatusOK, f)
}

func (a *App) deleteFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Film ID")
		return
	}

	f := models.Film{ID: id}

	if err := f.DeleteFilm(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

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

func (a *App) getTypes(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	cameras, err := models.GetTypes(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cameras)
}

func (a *App) createType(w http.ResponseWriter, r *http.Request) {
	var t models.Type
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}(r.Body)

	t.CreateType(a.DB)

	respondWithJSON(w, http.StatusCreated, t)
}

func (a *App) getType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	var c models.Type
	if !c.GetType(a.DB, id) {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) updateType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var c models.Type
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	c.ID = id
	c.UpdateType(a.DB)
	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Type ID")
		return
	}

	c := models.Type{ID: id}
	c.DeleteType(a.DB)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

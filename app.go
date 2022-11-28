package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Print(connectionString)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}
func (a *App) Run(addr string) {
	fmt.Println("Listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while returning response")
	}
}

func (a *App) getCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	c := camera{ID: id}
	if c.getCamera(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Camera not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}
	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) getCameras(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching cameras")
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	cameras, err := getCameras(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cameras)
}

func (a *App) createCamera(w http.ResponseWriter, r *http.Request) {
	var c camera
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	if err := c.createCamera(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) updateCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var c camera
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	c.ID = id

	if err := c.updateCamera(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
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

	c := camera{ID: id}
	if err := c.deleteCamera(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/cameras", a.getCameras).Methods("GET")
	a.Router.HandleFunc("/camera", a.createCamera).Methods("POST")
	a.Router.HandleFunc("/camera/{id:[0-9]+}", a.getCamera).Methods("GET")
	a.Router.HandleFunc("/camera/{id:[0-9]+}", a.updateCamera).Methods("PUT")
	a.Router.HandleFunc("/camera/{id:[0-9]+}", a.deleteCamera).Methods("DELETE")
}

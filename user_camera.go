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

func (a *App) getUserCameras(w http.ResponseWriter, r *http.Request) {
	token, err := a.GetTokenFromJWT(r)
	if err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 011)
		return
	}

	userToken, err := a.ReadJWTClaims(token)
	if err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 011)
		return
	}

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)

	var userCameraFilms []models.UserCamera
	if err = a.DB.Where("user_id = ?", user.ID).Find(&userCameraFilms).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, userCameraFilms)
}

func (a *App) createUserCamera(w http.ResponseWriter, r *http.Request) {
	token, err := a.GetTokenFromJWT(r)
	if err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 011)
		return
	}

	userToken, err := a.ReadJWTClaims(token)
	if err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 011)
		return
	}

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)

	var uc models.UserCamera
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&uc); err != nil {
		a.respondWithAnalogError(w, http.StatusBadRequest, 001)
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}(r.Body)

	var camera models.Camera
	if err = a.DB.First(&camera, uc.CameraID).Error; err != nil {
		a.respondWithAnalogError(w, http.StatusBadRequest, 013)
		return
	}

	uc.UserID = int(user.ID)
	if err = a.DB.Create(&uc).Error; err != nil {
		a.respondWithAnalogError(w, http.StatusInternalServerError, 014)
		return
	}

	respondWithJSON(w, http.StatusCreated, uc)
}

func (a *App) deleteUserCamera(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	token, err := a.GetTokenFromJWT(r)
	if err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 011)
		return
	}

	userToken, err := a.ReadJWTClaims(token)
	if err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 011)
		return
	}

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)

	var uc models.UserCamera
	if err := a.DB.First(&uc, id).Error; err != nil {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 017)
		return
	}

	if uc.UserID != int(user.ID) {
		a.respondWithAnalogError(w, http.StatusUnauthorized, 015)
		return
	}

	if err := a.DB.Delete(&uc).Error; err != nil {
		a.respondWithAnalogError(w, http.StatusInternalServerError, 016)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

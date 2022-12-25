package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/quentinchampenois/analog_api/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *App) getUserCameraFilms(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return a.JWTSecret, nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Your Token has been expired")
		return
	}

	var userToken struct {
		id     float64
		pseudo string
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == nil || claims["user_id"] == "" {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't contain required claims")
			return
		}
		if claims["pseudo"] == nil || claims["pseudo"] == "" {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't contain required claims")
			return
		}

		userToken.id = claims["user_id"].(float64)
		userToken.pseudo = claims["pseudo"].(string)
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)

	var uc models.UserCamera
	if err := a.DB.First(&uc, id).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if uc.UserID != int(user.ID) {
		respondWithError(w, http.StatusUnauthorized, "Not auhtorized to perform this action")
		return
	}

	var userCameraFilms []models.UserCameraFilm
	a.DB.Where("user_camera_id = ?", uc.ID).Find(&userCameraFilms)
	respondWithJSON(w, http.StatusOK, userCameraFilms)
}

func (a *App) createUserCameraFilms(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return a.JWTSecret, nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Your Token has been expired")
		return
	}

	var userToken struct {
		id     float64
		pseudo string
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == nil || claims["user_id"] == "" {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't contain required claims")
			return
		}
		if claims["pseudo"] == nil || claims["pseudo"] == "" {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't contain required claims")
			return
		}

		userToken.id = claims["user_id"].(float64)
		userToken.pseudo = claims["pseudo"].(string)
	}

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)

	var uc models.UserCameraFilm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&uc); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}(r.Body)

	var userCamera models.UserCamera
	if err := a.DB.First(&userCamera, uc.UserCameraID).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if userCamera.UserID != int(user.ID) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized action")
		return
	}

	uc.StartDate = time.Now()

	if err = a.DB.Create(&uc).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, uc)
}

func (a *App) rewindFilm(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return a.JWTSecret, nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Your Token has been expired")
		return
	}

	var userToken struct {
		id     float64
		pseudo string
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == nil || claims["user_id"] == "" {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't contain required claims")
			return
		}
		if claims["pseudo"] == nil || claims["pseudo"] == "" {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't contain required claims")
			return
		}

		userToken.id = claims["user_id"].(float64)
		userToken.pseudo = claims["pseudo"].(string)
	}

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var uc models.UserCameraFilm
	if err := a.DB.First(&uc, id).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var userCamera models.UserCamera
	if err := a.DB.First(&userCamera, uc.UserCameraID).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if userCamera.UserID != int(user.ID) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized action")
		return
	}

	uc.EndDate = time.Now();
	if err = a.DB.Save(&uc).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, uc)
}

package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/quentinchampenois/analog_api/models"
	"net/http"
	"strconv"
)

func (a *App) getUserCameraFilms(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

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
	fmt.Println(a.DB.Table("camera_films").Where("user_id = ?", user.ID).Where("camera_id = ?", id))
	var cameraFilms []models.CameraFilm
	a.DB.Where("user_id = ?", user.ID).Preload("User").Preload("Film").Preload("Camera").Find(&cameraFilms)
	respondWithJSON(w, http.StatusOK, cameraFilms)
}

/*
func (a *App) registerUserCamera(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	var c models.Camera
	if !c.GetCamera(a.DB, id) {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	if err := user.RegisterCamera(a.DB, &c); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "Successfully created !")
}

func (a *App) deleteUserCamera(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("An error occured on token parsing")
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
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	var c models.Camera
	if !c.GetCamera(a.DB, id) {
		respondWithError(w, http.StatusNotFound, "Camera not found")
		return
	}

	if err := user.DeleteCamera(a.DB, &c); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "Successfully deleted !")
}
*/

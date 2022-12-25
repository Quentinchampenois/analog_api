package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/quentinchampenois/analog_api/configs"
	"github.com/quentinchampenois/analog_api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router    *mux.Router
	DB        *gorm.DB
	Configs   configs.Config
	JWTSecret []byte
}

type Authentication struct {
	Pseudo   string `json:"pseudo"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Pseudo      string `json:"pseudo"`
	TokenString string `json:"token"`
}

func (a *App) Initialize() {
	a.Configs.Load()
	a.DB = a.Configs.Database.Connection()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run() {
	fmt.Printf("Listening on %v\n", a.Configs.Server.GetFullPath())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", a.Configs.Server.Port), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
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

func (a *App) signUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if user.Pseudo == "" || user.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Missing pseudo or password")
		return
	}

	user.Password, err = a.EncryptedPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	a.DB.Create(&user)
	respondWithJSON(w, http.StatusOK, user)
}

func (a *App) signIn(w http.ResponseWriter, r *http.Request) {
	var authdetails Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var authuser models.User
	a.DB.Where("pseudo = ?", authdetails.Pseudo).First(&authuser)
	if authuser.Pseudo == "" {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	if check := checkPasswordHash(authdetails.Password, authuser.Password); !check {
		respondWithError(w, http.StatusNotFound, "User pseudo or password is invalid")
		return
	}

	validToken, err := a.generateJWT(authuser.ID, authuser.Pseudo, authuser.Role)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to generate token")
		return
	}

	token := Token{
		Pseudo:      authuser.Pseudo,
		Role:        authuser.Role,
		TokenString: validToken,
	}
	respondWithJSON(w, http.StatusOK, token)
}

func (a *App) isAuthorized(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			respondWithError(w, http.StatusUnauthorized, "No Token Found")
			return
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

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return

			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}
		}
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
	})
}

func (a *App) generateJWT(userID uint, pseudo, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = userID
	claims["pseudo"] = pseudo
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	tokenString, err := token.SignedString(a.JWTSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *App) EncryptedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

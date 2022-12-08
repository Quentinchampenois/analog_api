package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/quentinchampenois/analog_api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type App struct {
	Router    *mux.Router
	DB        *gorm.DB
	JWTSecret []byte
}

type User struct {
	gorm.Model
	Pseudo   string `gorm:"unique;not null" json:"pseudo"`
	Password string `json:"password"`
	Role     string `json:"role"`
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

func (a *App) Initialize(host, port, user, password, dbname string) {
	dsn :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/paris", host, port, user, password, dbname)

	var err error
	a.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	if os.Getenv("JWT_SECRET_PASSWORD") == "" {
		log.Fatalln("You must define 'JWT_SECRET_PASSWORD' environment variable for JWT authentication system")
	}
	a.JWTSecret = []byte(os.Getenv("JWT_SECRET_PASSWORD"))
}
func (a *App) Run(addr string) {
	fmt.Println("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
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

func (a *App) getCameras(w http.ResponseWriter, r *http.Request) {
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
	var c Camera
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}(r.Body)

	if !c.createCamera(a.DB) {
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

	var c Camera
	if !c.getCamera(a.DB, id) {
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
	var c Camera
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	c.ID = id
	if err := c.updateCamera(a.DB); err != nil {
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

	c := Camera{ID: id}
	if err := c.deleteCamera(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

func (a *App) getTypes(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	cameras, err := getTypes(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cameras)
}

func (a *App) createType(w http.ResponseWriter, r *http.Request) {
	var c Type
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}(r.Body)

	c.createType(a.DB)

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) getType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid camera ID")
		return
	}

	var c Type
	if !c.getType(a.DB, id) {
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
	var c Type
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	c.ID = id
	c.updateType(a.DB)
	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) deleteType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Type ID")
		return
	}

	c := Type{ID: id}
	c.deleteType(a.DB)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})
}

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

func (a *App) signUp(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(user.Pseudo)
	if user.Pseudo == "" || user.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Missing pseudo or password")
		return
	}

	var dbuser User
	a.DB.Where("email = ?", user.Pseudo).First(&dbuser)

	//checks if email is already register or not
	if dbuser.Pseudo != "" {
		respondWithError(w, http.StatusNotFound, "Pseudo not found")
		return
	}

	user.Password, err = encryptedPassword(user.Password)
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

	var authuser User
	a.DB.Where("pseudo = ?", authdetails.Pseudo).First(&authuser)
	if authuser.Pseudo == "" {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	check := checkPasswordHash(authdetails.Password, authuser.Password)

	if !check {
		respondWithError(w, http.StatusNotFound, "User pseudo or password is invalid")
		return
	}

	validToken, err := a.generateJWT(authuser.Pseudo, authuser.Role)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to generate token")
		return
	}

	var token Token
	token.Pseudo = authuser.Pseudo
	token.Role = authuser.Role
	token.TokenString = validToken
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

func encryptedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *App) initializeRoutes() {
	cameraRouter := a.Router.PathPrefix("/camera").Subrouter()
	cameraRouter.HandleFunc("", a.getCameras).Methods("GET")
	cameraRouter.HandleFunc("", a.createCamera).Methods("POST")
	cameraRouter.HandleFunc("/", a.getCameras).Methods("GET")
	cameraRouter.HandleFunc("/", a.createCamera).Methods("POST")
	cameraRouter.HandleFunc("/{id:[0-9]+}", a.getCamera).Methods("GET")

	typeRouter := a.Router.PathPrefix("/type").Subrouter()
	typeRouter.HandleFunc("", a.getTypes).Methods("GET")
	typeRouter.HandleFunc("", a.createType).Methods("POST")
	typeRouter.HandleFunc("/", a.getTypes).Methods("GET")
	typeRouter.HandleFunc("/", a.createType).Methods("POST")
	typeRouter.HandleFunc("/{id:[0-9]+}", a.getType).Methods("GET")

	filmRouter := a.Router.PathPrefix("/film").Subrouter()
	filmRouter.HandleFunc("", a.getFilms).Methods("GET")
	filmRouter.HandleFunc("", a.createFilm).Methods("POST")
	filmRouter.HandleFunc("/", a.getFilms).Methods("GET")
	filmRouter.HandleFunc("/", a.createFilm).Methods("POST")
	filmRouter.HandleFunc("/{id:[0-9]+}", a.getFilm).Methods("GET")

	a.Router.HandleFunc("/signup", a.signUp).Methods("POST")
	a.Router.HandleFunc("/signin", a.signIn).Methods("POST")

	a.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})

	putProtectedRouter := a.Router.Methods(http.MethodPut).Subrouter()
	putProtectedRouter.HandleFunc("/camera/{id:[0-9]+}", a.updateCamera).Methods("PUT")
	putProtectedRouter.HandleFunc("/type/{id:[0-9]+}", a.updateType).Methods("PUT")
	putProtectedRouter.HandleFunc("/film/{id:[0-9]+}", a.updateFilm).Methods("PUT")

	deleteProtectedRouter := a.Router.Methods(http.MethodDelete).Subrouter()
	deleteProtectedRouter.HandleFunc("/camera/{id:[0-9]+}", a.deleteCamera).Methods("DELETE")
	deleteProtectedRouter.HandleFunc("/type/{id:[0-9]+}", a.deleteType).Methods("DELETE")
	deleteProtectedRouter.HandleFunc("/film/{id:[0-9]+}", a.deleteFilm).Methods("DELETE")

	putProtectedRouter.Use(a.isAuthorized)
	deleteProtectedRouter.Use(a.isAuthorized)
}

func (a *App) generateJWT(email, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	tokenString, err := token.SignedString(a.JWTSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

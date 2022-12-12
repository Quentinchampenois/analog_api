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
	"io"
	"log"
	"net/http"
	"strconv"
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

func (a *App) getCameras(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
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
			log.Fatalln(err)
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

func (a *App) getUserCameras(w http.ResponseWriter, r *http.Request) {
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

	var user models.User
	a.DB.Where("pseudo = ?", userToken.pseudo).Where("id = ?", userToken.id).First(&user)
	a.DB.Preload("Cameras").Find(&user)
	respondWithJSON(w, http.StatusOK, user.Cameras)
}

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

	var dbuser models.User
	a.DB.Where("pseudo = ?", user.Pseudo).First(&dbuser)

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

	var authuser models.User
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

	validToken, err := a.generateJWT(authuser.ID, authuser.Pseudo, authuser.Role)
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

	userRouter := a.Router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/cameras", a.getUserCameras).Methods("GET")
	userRouter.HandleFunc("/cameras/{id:[0-9]+}", a.registerUserCamera).Methods("POST")
	userRouter.Use(a.isAuthorized)
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

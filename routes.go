package main

import "net/http"

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

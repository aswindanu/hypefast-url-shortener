package app

import (
	"fmt"
	"log"
	"net/http"

	// handler

	"hypefast-url-shortener/app/handler/urlshortener"

	// model
	"hypefast-url-shortener/app/model"

	// config
	"hypefast-url-shortener/config"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	var err error
	var db *gorm.DB
	var dbURI string

	// MYSQL
	if config.DB.Dialect == "mysql" {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
			config.DB.Username,
			config.DB.Password,
			config.DB.Host,
			config.DB.Port,
			config.DB.Name,
			config.DB.Charset,
		)
	}
	// POSTGRES
	if config.DB.Dialect == "postgres" {
		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			config.DB.Host,
			config.DB.Username,
			config.DB.Password,
			config.DB.Name,
			config.DB.Port,
		)
	}

	db, err = gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		fmt.Println("dbURI: ", dbURI)
		log.Fatal("Could not connect database:" + err.Error())
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the urlshortener
	a.Get("/api/v1/url-shortener", a.handleRequest(urlshortener.GetAllUrls))
	a.Post("/api/v1/url-shortener", a.handleRequest(urlshortener.CreateUrl))
	a.Get("/api/v1/url-shortener/redirect", a.handleRequest(urlshortener.RedirectUrl))
	a.Get("/api/v1/url-shortener/{uid}", a.handleRequest(urlshortener.GetUrl))
	a.Post("/api/v1/url-shortener/{uid}/activate", a.handleRequest(urlshortener.ActivateUrl))
	a.Post("/api/v1/url-shortener/{uid}/deactivate", a.handleRequest(urlshortener.DeactivateUrl))
	// a.Delete("/api/v1/url-shortener/{uid}/deactivate", a.handleRequest(urlshortener.DeactivateProject)) NOTE: User should not delete data (changed to activate/deactivate)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

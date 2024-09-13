package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/user-manager/app"
	"github.com/user-manager/handlers"
)

func NewRouter(app *app.App) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", handlers.GetUsers(app)).Methods("GET")
	r.HandleFunc("/users", handlers.AddUsers(app)).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.DeleteUsers(app)).Methods("DELETE")
	r.HandleFunc("/users/{id}", handlers.EditUsers(app)).Methods("PATCH")
	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return r
}

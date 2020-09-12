package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilikerice123/puzzle/game"
	"github.com/ilikerice123/puzzle/store"
)

// RegisterUsersRoutes registers /api/images routers
func RegisterUsersRoutes(r *mux.Router) {
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("", CreateUser).Methods("POST")
	usersRouter.HandleFunc("/", CreateUser).Methods("POST")
	usersRouter.HandleFunc("/{id}", GetUser).Methods("GET")
	usersRouter.HandleFunc("/{id}/", GetUser).Methods("GET")
	usersRouter.HandleFunc("/auth", AuthUser).Methods("GET")
	usersRouter.HandleFunc("/auth/", AuthUser).Methods("GET")
}

// GetUser gets a user given an id
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if user := game.GlobalUserPool.GetUser(id); user != nil {
		WriteSuccess(w, user)
		return
	}
	WriteError(w, 404, map[string]string{"error": "player not found"})
}

// AuthUser returns the uuid of a user given a username and password
func AuthUser(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok {
		WriteError(w, 401, map[string]string{"error": "not authenticated properly"})
	}
	WriteSuccess(w, map[string]string{"user": user, "pass": password})
}

// CreateUser creates a user given a string
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInfo map[string]string
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": err.Error()})
		return
	}
	name := userInfo["name"]
	password := userInfo["password"]
	if name == "" || password == "" {
		WriteError(w, 422, map[string]string{"error": "must provide name and password"})
		return
	}
	user := store.NewUser(name, password)
	game.GlobalUserPool.AddUser(user)
	WriteSuccess(w, user)
}

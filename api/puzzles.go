package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilikerice123/puzzle/game"
)

// RegisterPuzzlesRoutes registers /api/images routers
func RegisterPuzzlesRoutes(r *mux.Router) {
	usersRouter := r.PathPrefix("/puzzles").Subrouter()
	usersRouter.HandleFunc("/{id}", GetPuzzle).Methods("GET")
	usersRouter.HandleFunc("/{id}/", GetPuzzle).Methods("GET")
	usersRouter.HandleFunc("/{id}", CreatePuzzle).Methods("GET")
	usersRouter.HandleFunc("/{id}/", CreatePuzzle).Methods("GET")
}

// GetPuzzle gets current puzzle's state
func GetPuzzle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if user := game.GlobalUserPool.GetUser(id); user != nil {
		WriteSuccess(w, user)
		return
	}
	WriteError(w, 404, map[string]string{"error": "player not found"})
}

// CreatePuzzle creates a puzzle given a size
func CreatePuzzle(w http.ResponseWriter, r *http.Request) {
	var userInfo map[string]string
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": err.Error()})
		return
	}
	name := userInfo["name"]
	if name == "" {
		WriteError(w, 422, map[string]string{"error": "must provide name"})
		return
	}
	user := game.NewUser(name)
	game.GlobalUserPool.AddUser(user)
	WriteSuccess(w, user)
}

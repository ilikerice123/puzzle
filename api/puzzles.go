package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilikerice123/puzzle/fs"
	"github.com/ilikerice123/puzzle/game"
)

// RegisterPuzzlesRoutes registers /api/images routers
func RegisterPuzzlesRoutes(r *mux.Router) {
	usersRouter := r.PathPrefix("/puzzles").Subrouter()
	usersRouter.HandleFunc("/{id}", GetPuzzle).Methods("GET")
	usersRouter.HandleFunc("/{id}/", GetPuzzle).Methods("GET")
	usersRouter.HandleFunc("/{id}", CreatePuzzle).Methods("POST")
	usersRouter.HandleFunc("/{id}/", CreatePuzzle).Methods("POST")
}

// GetPuzzle gets current puzzle's state
func GetPuzzle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	puzzle := game.GlobalPuzzlePool.GetPuzzle(id)
	if puzzle != nil {
		WriteSuccess(w, puzzle)
		return
	}
	WriteError(w, 404, map[string]string{"error": "puzzle not found"})
}

// CreatePuzzle creates a puzzle given a size
func CreatePuzzle(w http.ResponseWriter, r *http.Request) {
	var userInfo map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": err.Error()})
		return
	}
	id, ok := userInfo["id"].(string)
	pictureFile := "images/" + id + "/original.jpeg"
	if !ok || !fs.DirExists(pictureFile) {
		WriteError(w, 422, map[string]string{"error": "invalid id provided"})
	}
	ySize, oky := userInfo["y_size"].(int)
	xSize, okx := userInfo["x_size"].(int)
	if !oky || !okx || ySize <= 0 || xSize <= 0 {
		WriteError(w, 422, map[string]string{"error": "invalid x_size and y_size provided"})
		return
	}

	puzzle := game.NewLivePuzzle(pictureFile, ySize, xSize, game.GlobalUserPool)
	game.GlobalPuzzlePool.AddPuzzle(puzzle)
	if puzzle == nil {
		WriteError(w, 500, map[string]string{"error": "error creating puzzle"})
		return
	}
	WriteSuccess(w, map[string]string{"id": id})
}

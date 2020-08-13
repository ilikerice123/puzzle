package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ilikerice123/puzzle/fs"
	"github.com/ilikerice123/puzzle/game"
)

// RegisterPuzzlesRoutes registers /api/images routers
func RegisterPuzzlesRoutes(r *mux.Router) {
	puzzlesRouter := r.PathPrefix("/puzzles").Subrouter()
	puzzlesRouter.HandleFunc("/{id}/ws", UpgradePuzzle)
	puzzlesRouter.HandleFunc("/{id}", GetPuzzle).Methods("GET")
	puzzlesRouter.HandleFunc("/{id}/", GetPuzzle).Methods("GET")
	puzzlesRouter.HandleFunc("/{id}", CreatePuzzle).Methods("POST")
	puzzlesRouter.HandleFunc("/{id}/", CreatePuzzle).Methods("POST")
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
	var userInfo map[string]int
	id := mux.Vars(r)["id"]
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": err.Error()})
		return
	}

	pictureFile := "images/" + id + "/original.jpeg"
	if !fs.DirExists(pictureFile) {
		WriteError(w, 422, map[string]string{"error": "invalid id provided"})
	}

	ySize := userInfo["y_size"]
	xSize := userInfo["x_size"]
	if ySize <= 0 || xSize <= 0 {
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

// UpgradePuzzle creates puzzle
func UpgradePuzzle(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	puzzle := game.GlobalPuzzlePool.GetPuzzle(id)
	if puzzle != nil {
		WriteError(w, 404, map[string]string{"error": "puzzle does not exist"})
	}
	log.Println("trying to connect and upgrade!")
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(id))

}

func setUpConnection(c *websocket.Conn, p game.LivePuzzleBase) {
	for {
		messageType, p, err := c.ReadMessage()
		if messageType == websocket.BinaryMessage {
			log.Println("binary message!")
		} else {
			log.Println("text message!")
		}
		if err != nil {
			log.Println(err)
			return
		}
		if err := c.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

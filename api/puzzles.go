package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	ySize := userInfo["ySize"]
	xSize := userInfo["xSize"]
	if ySize <= 0 || xSize <= 0 {
		WriteError(w, 422, map[string]string{"error": "invalid xSize and ySize provided"})
		return
	}
	puzzle := game.NewLivePuzzle(id, pictureFile, ySize, xSize, game.GlobalUserPool)
	puzzle.Start()
	game.GlobalPuzzlePool.AddPuzzle(puzzle)
	if puzzle == nil {
		WriteError(w, 500, map[string]string{"error": "error creating puzzle"})
		return
	}
	WriteSuccess(w, map[string]string{"id": id})
}

// UpgradePuzzle creates puzzle socket
func UpgradePuzzle(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	if userID == "" {
		WriteError(w, 404, map[string]string{"error": "user parameter not supplied"})
		return
	}

	id := mux.Vars(r)["id"]
	puzzle := game.GlobalPuzzlePool.GetPuzzle(id)
	if puzzle != nil {
		WriteError(w, 404, map[string]string{"error": "puzzle does not exist"})
		return
	}

	log.Println("trying to connect and upgrade!")
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	go setupConnection(conn, puzzle, userID)
}

// setupConnection connects all the pipelines and channels together
func setupConnection(c *websocket.Conn, p game.LivePuzzleBase, userID string) {
	defer c.Close()
	// set up heartbeats
	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// pushing updates path
	p.AddCallback(func(u *game.Update) {
		serializedUpdate, err := json.Marshal(u)
		if err != nil {
			return
		}
		c.WriteMessage(websocket.TextMessage, serializedUpdate)
	})

	p.AddRequest(&game.Request{Action: game.JOIN, UserID: userID})
	// wire up connections first, then send join message, so we also get connected message
	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil || msgType != websocket.TextMessage {
			log.Println(err)
			return
		}
		var r game.Request
		if err := json.Unmarshal(msg, &r); err != nil {
			continue
		}
		p.AddRequest(&r)
	}
}

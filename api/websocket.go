package api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketUpgrader is the global websocket upgrader object
var WebsocketUpgrader websocket.Upgrader

// InitUpgrader used to initialize the global websocket upgrader object
func InitUpgrader() {
	WebsocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 4096}
	WebsocketUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}

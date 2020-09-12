package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/ilikerice123/puzzle/api"
	"github.com/ilikerice123/puzzle/fs"
	"github.com/ilikerice123/puzzle/game"
	"github.com/ilikerice123/puzzle/store"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("Hello!!! serving traffic")
	// picture.SliceImage("images/6bfdfa70-d02e-4ba7-a65b-b35627b22212/original.jpeg", 3, 4)

	// make directory to store images
	if !fs.DirExists("images") {
		err := os.Mkdir("images", 0666)
		if err != nil {
			log.Fatalf("unable to create images directory to store images")
		}
	}

	// init global pools and websocket upgrader
	game.InitUserPool()
	game.InitPuzzlePool()
	api.InitUpgrader()
	store.InitStore()

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	api.RegisterImagesRoutes(apiRouter)
	api.RegisterUsersRoutes(apiRouter)
	api.RegisterPuzzlesRoutes(apiRouter)
	api.RegisterFrontEnd(r)

	router := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://ec2-54-245-184-188.us-west-2.compute.amazonaws.com"},
		AllowCredentials: true}).Handler(r)

	server := &http.Server{
		Handler: router,
		Addr:    ":80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 24 * time.Hour,
		ReadTimeout:  24 * time.Hour,
	}
	log.Fatal(server.ListenAndServe())
	wait := make(chan int, 1)
	<-wait
}

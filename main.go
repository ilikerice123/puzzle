package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/ilikerice123/puzzle/api"
)

func main() {
	fmt.Println("Hello!!! serving traffic")
	// picture.SliceImage("images/6bfdfa70-d02e-4ba7-a65b-b35627b22212/original.jpeg", 3, 4)

	//make directory to store images
	_, err := os.Stat("images")
	if os.IsNotExist(err) {
		err = os.Mkdir("images", 0666)
	}

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	imagesRouter := apiRouter.PathPrefix("/images").Subrouter()
	api.RegisterImageRoutes(imagesRouter)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 24 * time.Hour,
		ReadTimeout:  24 * time.Hour,
	}
	log.Fatal(srv.ListenAndServe())
	wait := make(chan int, 1)
	<-wait
}

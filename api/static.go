package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterFrontEnd registers the react front end for "/"
func RegisterFrontEnd(r *mux.Router) {
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "client/build/index.html")
	})
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("client/build/static")))
	r.PathPrefix("/static/").Handler(staticHandler)
}

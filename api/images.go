package api

import (
	"image"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ilikerice123/puzzle/fs"
	"github.com/ilikerice123/puzzle/picture"
)

// RegisterImagesRoutes registers /api/images routers
func RegisterImagesRoutes(r *mux.Router) {
	imagesRouter := r.PathPrefix("/images").Subrouter()
	imagesRouter.Methods("GET").Handler(
		http.StripPrefix("/api/", http.FileServer(http.Dir(""))))
	imagesRouter.HandleFunc("", UploadImage).Methods("POST")
	imagesRouter.HandleFunc("/", UploadImage).Methods("POST")
}

// UploadImage uploads an image to a directory, and creates a preview
func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // 32MB limit
	file, _, err := r.FormFile("image")
	if err != nil {
		// file error
		WriteError(w, 500, map[string]string{"error": err.Error()})
		return
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": "error decoding image"})
		return
	}
	preview := picture.DownsizeImage(img)

	uuid := uuid.New().String()
	if fs.DirExists("images/" + uuid) {
		WriteError(w, 500, map[string]string{"error": "directory exists, probably uuid collision"})
		return
	}

	os.Mkdir("images/"+uuid, 0666)
	err = fs.SaveImage("images/"+uuid+"/original.jpeg", img)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": err.Error()})
		return
	}
	err = fs.SaveImage("images/"+uuid+"/preview.jpeg", preview)
	if err != nil {
		WriteError(w, 422, map[string]string{"error": err.Error()})
		return
	}
	WriteSuccess(w, map[string]string{"uuid": uuid})
}

package methuforecasttimelapse

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func gifHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)
	w.Header().Set("content-type", "image/gif")
	http.ServeFile(w, r, "gifs/anim.gif")
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)
	w.Header().Set("content-type", "text/html")
	t, err := template.ParseFiles("templates/gallery.template")
	if err != nil {
		log.Println(err)
		return
	}
	images, err := ListImageFiles("images")
	t.Execute(w, images)
}

func initHandlers() {
	staticServer := http.FileServer(http.Dir("static"))
	imageServer := http.FileServer(http.Dir("images"))
	http.HandleFunc("/gif", gifHandler)
	http.HandleFunc("/gallery", galleryHandler)
	http.Handle("/images/", http.StripPrefix("/images/", imageServer))
	http.Handle("/", staticServer)
}

func StartServer(localAddress string, port int) {
	initHandlers()

	serverAddress := fmt.Sprintf("%v:%v", localAddress, port)
	log.Println("Listening on", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Println("Problem with the server start", err)
	}
}

// For the Google App Engine initialization
func init() {
	initHandlers()
}

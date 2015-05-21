package main

import (
	"fmt"
	"log"
	"net/http"
)

func gifHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)
	w.Header().Set("content-type", "image/gif")
	http.ServeFile(w, r, "gifs/anim.gif")
}

func StartServer(localAddress string, port int) {
	staticServer := http.FileServer(http.Dir("static"))
	http.HandleFunc("/gif", gifHandler)
	http.Handle("/", staticServer)

	serverAddress := fmt.Sprintf("%v:%v", localAddress, port)
	log.Println("Listening on", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Println("Problem with the server start", err)
	}
}

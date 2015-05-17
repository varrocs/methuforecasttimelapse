package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("mainpage.html")
	if err != nil {
		t.Execute(w, nil)
	}
}

func handleGif(w http.ResponseWriter, r *http.Request) {
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
}

func handleLog(w http.ResponseWriter, r *http.Request) {
}

func StartServer() {
	port := flag.Int("port", 8080, "Web server listening port")
	address := flag.String("address", "0.0.0.0", "Local address to bind to")

	flag.Parse()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/gif", handleGif)
	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/log", handleLog)

	err := http.ListenAndServe(fmt.Sprintf("%v:%v", *address, *port), nil)
	if err != nil {
		log.Println("Problem with the server start ", err)
	}
}

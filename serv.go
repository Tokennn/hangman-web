package main

import (
	"log"
	"net/http"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./home.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/ingame", Inside)
	http.HandleFunc("/vidhome", Video)
	http.HandleFunc("/endpage", Exit)

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func Inside(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./ingame.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, nil)
}

func Video(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./vidhome.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, nil)
}

func Exit(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./endpage.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, nil)
}

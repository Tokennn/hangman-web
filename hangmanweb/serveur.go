package hangmanweb

import (
	"log"
	"net/http"
	"text/template"

	"github.com/Tokennn/hangman"
)

type HangmanData struct {
	Word    string
	Display string
}

func Home(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	template, err := template.ParseFiles("./home.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, data)
}

func Serveur() {

	data := &HangmanData{}
	data.Word = hangman.Randomly()
	data.Display = hangman.Displaywords(data.Word)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { Home(w, r, data) })
	http.HandleFunc("/ingame", func(w http.ResponseWriter, r *http.Request) { Inside(w, r, data) })
	http.HandleFunc("/vidhome", func(w http.ResponseWriter, r *http.Request) { Video(w, r, data) })
	http.HandleFunc("/endpage", func(w http.ResponseWriter, r *http.Request) { Exit(w, r, data) })

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}

func Inside(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	template, err := template.ParseFiles("./ingame.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, data)
}

func Video(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	template, err := template.ParseFiles("./vidhome.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, data)
}

func Exit(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	template, err := template.ParseFiles("./endpage.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, data)
}

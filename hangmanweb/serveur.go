package hangmanweb

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/Tokennn/hangman"
)

type HangmanData struct {
	Word      string
	Display   string
	Letters   string
	Life      int
	Useletter []string
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
	fmt.Println(data.Word)
	fmt.Println(data.Letters)
	data.Display = hangman.Displaywords(data.Word)
	data.Life = 10

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { Home(w, r, data) })
	http.HandleFunc("/ingame", func(w http.ResponseWriter, r *http.Request) { Inside(w, r, data) })
	http.HandleFunc("/endpage", func(w http.ResponseWriter, r *http.Request) { Exit(w, r, data) })
	http.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) { Restart(w, r, data) })
	http.HandleFunc("/putletter", func(w http.ResponseWriter, r *http.Request) { PutLetter(w, r, data) })

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}

func Inside(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	template, err := template.ParseFiles("./ingame.html")
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

func Restart(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	data.Word = hangman.Randomly()
	data.Display = hangman.Displaywords(data.Word)
	Inside(w, r, data)
	template, err := template.ParseFiles("./ingame.html")
	if err != nil {
		log.Fatal(err)
	}

	template.Execute(w, data)
}

func PutLetter(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	Letters := r.FormValue("letter")
	data.Letters = Letters

	if Letters == "" {
		Inside(w, r, data)
		return
	}

	if len(data.Letters) > 1 {
		if data.Letters == data.Word {
			http.Redirect(w, r, "/endpage", http.StatusSeeOther)
			return
		} else {
			data.Life -= 2
		}
	} else {
		used := false
		find := false
		for _, letter := range data.Useletter {
			if letter == Letters {
				used = true
				break
			}
		}
		if used {
			http.Redirect(w, r, "/ingame", http.StatusSeeOther)
			return
		} else {
			for i, char := range data.Word {
				if string(char) == Letters {
					find = true
					newWord := ""
					for x := range data.Display {
						if x == i {
							newWord += string(char)
						} else {
							newWord += string(data.Display[x])
						}
					}
					data.Display = newWord
				}
			}
			if !find {
				data.Life -= 1
			}
			data.Useletter = append(data.Useletter, Letters)
		}
	}
	if data.Life <= 0 {
		http.Redirect(w, r, "/endpage", http.StatusSeeOther)
		return
	}
	if data.Display == data.Word {
		http.Redirect(w, r, "/endpage", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/ingame", http.StatusSeeOther)
}

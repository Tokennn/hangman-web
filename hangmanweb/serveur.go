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
	Message   string
	Message2  string
	Username  string
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
	data.Life = 10

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { Home(w, r, data) })
	http.HandleFunc("/ingame", func(w http.ResponseWriter, r *http.Request) { Inside(w, r, data) })
	http.HandleFunc("/endpage", func(w http.ResponseWriter, r *http.Request) { Endpage(w, r, data) })
	http.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) { Restart(w, r, data) })
	http.HandleFunc("/putletter", func(w http.ResponseWriter, r *http.Request) { PutLetter(w, r, data) })
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) { Win(w, r, data) })
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { Register(w, r, data) })

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

func Endpage(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	data.Message2 = " The word was " + data.Word
	template, err := template.ParseFiles("./endpage.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}

func Win(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	template, err := template.ParseFiles("./winpage.html")
	if err != nil {
		log.Fatal(err)

	}
	template.Execute(w, data)
}

func Restart(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	data.Word = hangman.Randomly()
	data.Display = hangman.Displaywords(data.Word)
	data.Life = 10
	data.Useletter = []string{}
	http.Redirect(w, r, "/ingame", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	data.Username = r.FormValue("username")

	if r.Method == "GET" {

		http.ServeFile(w, r, "./static/register.html")

	} else if r.Method == "POST" {
		r.ParseForm()

		http.SetCookie(w, &http.Cookie{
			Name:  "username",
			Value: data.Username,
		})

		http.Redirect(w, r, "/ingame", http.StatusSeeOther)
	}
}

func PutLetter(w http.ResponseWriter, r *http.Request, data *HangmanData) {
	Letters := r.FormValue("letter")
	data.Letters = Letters
	data.Message = ""

	if Letters == "" {
		Inside(w, r, data)
		return
	}

	if len(data.Letters) > 1 {
		if data.Letters == data.Word {
			http.Redirect(w, r, "/win", http.StatusSeeOther)
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
			data.Message = "🚧 You already used this letter 🚧"
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
		http.Redirect(w, r, "/win", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/ingame", http.StatusSeeOther)
	}
}

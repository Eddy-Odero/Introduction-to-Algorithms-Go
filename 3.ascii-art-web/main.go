package main

import (
	"html/template"
	"net/http"

	"ascii/ascii"
)

type PageData struct {
	Text   string
	Banner string
	Result string
	Error  string
}

var tmpl = template.Must(template.ParseFiles("template/index.html"))

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := PageData{
		Banner: "standard",
	}

	tmpl.Execute(w, data)
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	data := PageData{
		Text:   text,
		Banner: banner,
	}

	if text == "" {
		data.Error = "Please enter some text"
		tmpl.Execute(w, data)
		return
	}

	lines, err := ascii.LoadBanner("banners/" + banner + ".txt")
	if err != nil {
		data.Error = "Invalid banner file"
		tmpl.Execute(w, data)
		return
	}

	font := ascii.BuildFont(lines)
	data.Result = ascii.Generate(text, font)

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ascii-art", asciiHandler)

	http.ListenAndServe(":8080", nil)
}
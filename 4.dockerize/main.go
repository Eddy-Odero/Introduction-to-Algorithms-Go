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

	if err := tmpl.Execute(w, data); err != nil {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)//500 error
	return
}

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
	http.Error(w, "Bad Request: empty input", http.StatusBadRequest)//400 bad request
	return
}
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		data.Error = "Invalid banner selected"
		tmpl.Execute(w, data)
		return
	}

	lines, err := ascii.LoadBanner("banners/" + banner + ".txt")
if err != nil {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)// 500 error
	return
}

	font := ascii.BuildFont(lines)
	data.Result = ascii.Generate(text, font)

	tmpl.Execute(w, data)
}
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {//404 error
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	home(w, r)
})
	http.HandleFunc("/ascii-art", asciiHandler)

	http.ListenAndServe(":8080", nil)
}
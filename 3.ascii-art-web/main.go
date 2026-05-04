package main

import (
	"fmt"
	"net/http"
)
func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == ""{
		name = "Eddy"
	}
	fmt.Fprintf(w, "Hello %s", name)
}
func main(){
	http.HandleFunc("/",hello)
	http.ListenAndServe(":8080", nil)

}
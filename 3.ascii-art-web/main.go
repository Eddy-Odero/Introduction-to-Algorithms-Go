package main

import (
	"fmt"
	"net/http"
)
func home(w http.ResponseWriter, r * http.Request){
		http.ServeFile(w,r, "template/index.html")
	}
	func submit(w http.ResponseWriter, r * http.Request){
		if r.Method !=  http.MethodPost{
			http.Error(w , "Invalid Request", http.StatusMethodNotAllowed)
			return 
		}
		name := r.FormValue("name")
		fmt.Fprintf(w, "Hello :%s", name)
	}
func main(){
	http.HandleFunc("/", home)
	http.HandleFunc("/submit",submit)
	http.ListenAndServe(":8080", nil)
}
package main

import (
	"html/template"
	"net/http"
)

func main() {

	userTmpl := template.Must(template.ParseFiles("./template/user.html"))
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		var username string
		_, exists := r.Header["X-Blazewall-Session"]
		if exists {
			username = r.Header["X-Blazewall-Session"][0]
		}
		data := struct {
			Username string
		}{
			Username: username,
		}
		userTmpl.Execute(w, data)
	})
	indexTmlp := template.Must(template.ParseFiles("./template/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexTmlp.Execute(w, nil)
	})
	http.ListenAndServe(":8080", nil)
}

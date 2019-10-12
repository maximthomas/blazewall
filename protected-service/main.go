package main

import (
	"html/template"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	t := template.Must(template.ParseGlob("./template/*.*"))
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
		t.ExecuteTemplate(w, "user.html", data)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "index.html", nil)
	})
	http.ListenAndServe(":8080", nil)
}

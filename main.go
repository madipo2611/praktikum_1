package main

import (
	"html/template"
	"log"
	"net/http"
)

type Users struct {
	login    string
	password string
}

var users = map[int][]Users{}

func main() {
	users[1] = append(users[1], Users{login: "admin", password: "admin"})
	users[2] = append(users[2], Users{login: "admin2", password: "admin2"})
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/create_card", create_card)

	fileServer := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ts, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	if r.Method != http.MethodPost {
		ts.Execute(w, nil)
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if (username == "admin" && password == "admin") || (username == "admin2" && password == "admin2") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/register.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func create_card(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/create_card.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

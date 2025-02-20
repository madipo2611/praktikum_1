package main

import (
	"context"
	"database/sql"
	"html/template"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
)

type Users struct {
	ID       int
	Login    string
	Password string
}

type Cards struct {
	ID       int
	Title    string
	Subtitle string
	Text     string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "jurry")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/home", home)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/create_card", create_card)

	fileServer := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	err = http.ListenAndServe(":8080", mux)
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
		return
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
		return
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
		return
	}

	if r.Method != http.MethodPost {
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("Попытка входа: Логин: %s, Пароль: %s", username, password)

	var user Users // Объявляем переменную user для хранения данных пользователя

	err = db.QueryRowContext(r.Context(), "SELECT id FROM users WHERE Login = ? AND Password = ?", username, password).Scan(&user.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Пользователь не найден:", username)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			log.Println("Ошибка при выполнении запроса:", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	// Если пользователь найден и данные верные, перенаправляем на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func register(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/register.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if r.Method != http.MethodPost {
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	var exists bool
	err = db.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE Login = ?)", username).Scan(&exists)
	if err != nil {
		return
	}

	if exists {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Если логин не существует, добавляем нового пользователя
	_, err = db.ExecContext(context.Background(), "INSERT INTO users (Login, Password) VALUES (?, ?)", username, password)
	if err != nil {
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func create_card(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/create_card.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if r.Method != http.MethodPost {
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	title := r.FormValue("title")
	subtitle := r.FormValue("subtitle")
	text := r.FormValue("text")

	_, err = db.ExecContext(context.Background(), "INSERT INTO cards (Title, Subtitle, Text) VALUES (?, ?, ?)", title, subtitle, text)
	if err != nil {
		log.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

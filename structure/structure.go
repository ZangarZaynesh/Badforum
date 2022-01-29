package structure

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Error struct {
	Str    string
	number int
}

type HandlerArt struct {
	Art *sql.DB
}

type Posts struct {
	Id_users, Id_posts, User_login, User_password, User_email, Post_date, Post_user_id, Post_category_id, Post string
}

func Err(Str string, Status int, w http.ResponseWriter, r *http.Request) {

	Info := Error{Str, Status}
	val, err := template.ParseFiles("templates/error.html")

	if err != nil {
		log.Println("Error when parsing a template: %s", err)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(Status)
	err = val.ExecuteTemplate(w, "error.html", Info)
	if err != nil {
		log.Println("Error when parsing a template: %s", err)
		fmt.Fprintf(w, err.Error())
		return
	}
}

func (h *HandlerArt) Index(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if r.Method != "GET" {
		Err("405 Method Not Allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	var Post []Posts
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	rows, err := h.Art.Query("SELECT * FROM users INNER JOIN posts ON users.id = posts.user_id")
	if err != nil {
		log.Fatal(err)
	}
	var id_users, id_posts, user_login, user_password, user_email, post_date, post_user_id, post_category_id, post string
	for rows.Next() {
		err = rows.Scan(&id_users, &user_login, &user_password, &user_email, &id_posts, &post_date, &post_user_id, &post_category_id, &post)
		Post = append(Post, Posts{Id_users: id_users,
			Id_posts:         id_posts,
			User_login:       user_login,
			User_password:    user_password,
			User_email:       user_email,
			Post_date:        post_date,
			Post_user_id:     post_user_id,
			Post_category_id: post_category_id,
			Post:             post})
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(Post.Id_users, Post.Id_posts, Post.User_login, Post.User_password, Post.User_email, Post.Post_date, Post.Post_user_id, Post.Post_category_id, Post.Post)
	}
	defer rows.Close()

	err = tmpl.ExecuteTemplate(w, "index.html", Post) // h.Art.Page - необходимая структура/подструктура
	if err != nil {
		log.Println("Error when parsing a template:", err)
		fmt.Fprintf(w, err.Error())
		return
	}
}

func (h *HandlerArt) Registration(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/registration" {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if r.Method != "GET" {
		Err("405 Method Not Allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		fmt.Println(err)
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	err = tmpl.ExecuteTemplate(w, "registration.html", h.Art) // h.Art.Page - необходимая структура/подструктура
	if err != nil {
		log.Println("Error when parsing a template:", err)
		fmt.Fprintf(w, err.Error())
		return
	}
}

func (h *HandlerArt) Created(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/registration/created" {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if r.Method != "POST" {
		Err("405 Method Not Allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	// login := r.FormValue("login")
	// password := r.FormValue("password")
	// email := r.FormValue("email")

	tmpl, err := template.ParseFiles("templates/created.html")
	if err != nil {
		fmt.Println(err)
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	err = tmpl.ExecuteTemplate(w, "created.html", h.Art) // h.Art.Page - необходимая структура/подструктура
	if err != nil {
		log.Println("Error when parsing a template:", err)
		fmt.Fprintf(w, err.Error())
		return
	}
}

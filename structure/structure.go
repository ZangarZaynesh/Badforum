package structure

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Str    string
	number int
}

type HandlerArt struct {
	Art *sql.DB
}

type SendStruct struct {
	SuccessFull string
}

type Posts struct {
	Id_users, Id_posts, User_login, User_password, User_email, Post_date, Post_user_id, Post_category_id, Post string
}

func Err(Str string, Status int, w http.ResponseWriter, r *http.Request) {

	Info := Error{Str, Status}
	w.WriteHeader(Status)
	ExecTemp("templates/error.html", "error.html", Info, w, r)
}

func (h *HandlerArt) Index(w http.ResponseWriter, r *http.Request) {

	if !PathMethod("/", "GET", w, r) {
		return
	}

	rows, err := h.Art.Query("SELECT * FROM users INNER JOIN posts ON users.id = posts.user_id;")
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	var Post []Posts
	var id_users, id_posts, user_login, user_password, user_email, post_date, post_user_id, post_category_id, post string
	for rows.Next() {
		err = rows.Scan(&id_users, &user_login, &user_password, &user_email, &id_posts, &post_date, &post_user_id, &post_category_id, &post)
		if err != nil {
			log.Fatal(err)
		}
		Post = append(Post, Posts{Id_users: id_users,
			Id_posts:         id_posts,
			User_login:       user_login,
			User_password:    user_password,
			User_email:       user_email,
			Post_date:        post_date,
			Post_user_id:     post_user_id,
			Post_category_id: post_category_id,
			Post:             post})
	}

	defer rows.Close()
	ExecTemp("templates/index.html", "index.html", Post, w, r)
}

func (h *HandlerArt) Registration(w http.ResponseWriter, r *http.Request) {

	if !PathMethod("/registration", "GET", w, r) {
		return
	}
	ExecTemp("templates/registration.html", "registration.html", h.Art, w, r)
}

func (h *HandlerArt) Created(w http.ResponseWriter, r *http.Request) {

	if !PathMethod("/registration/created", "POST", w, r) {
		return
	}

	login := r.FormValue("login")
	if login == "" {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if !test("login", login, h, w, r) {
		return
	}

	password := r.FormValue("password")
	if password == "" {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	email := r.FormValue("email")
	if !isEmailValid(email) {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if !test("email", email, h, w, r) {
		return
	}

	GenPassword := []byte(password)
	if !GeneratePass(&GenPassword, w, r) {
		return
	}

	_, err := h.Art.Exec("INSERT INTO users (login, password, email) VALUES ( ?, ?, ?);", login, GenPassword, email)
	if err != nil {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	ExecTemp("templates/created.html", "created.html", &SendStruct{SuccessFull: "succesfull"}, w, r)
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func PathMethod(Path, Method string, w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != Path {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return false
	}

	if r.Method != Method {
		Err("405 Method Not Allowed", http.StatusMethodNotAllowed, w, r)
		return false
	}
	return true
}

func test(NameColumn, ValueColumn string, h *HandlerArt, w http.ResponseWriter, r *http.Request) bool {
	row := h.Art.QueryRow("SELECT ? FROM users WHERE users."+NameColumn+"= ?;", NameColumn, ValueColumn)

	err1 := row.Scan()
	if err1 == nil && err1 != sql.ErrNoRows {
		Err("400 Bad Request", 400, w, r)
		return false
	}

	return true
}

func GeneratePass(password *[]byte, w http.ResponseWriter, r *http.Request) bool {

	var err1 error
	*password, err1 = bcrypt.GenerateFromPassword(*password, 8)

	if err1 != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return false
	}
	return true
}

func ExecTemp(PathHTML, NameHTML string, Struct interface{}, w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(PathHTML)
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	err = tmpl.ExecuteTemplate(w, NameHTML, Struct)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

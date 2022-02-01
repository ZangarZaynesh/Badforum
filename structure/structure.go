package structure

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/satori/uuid"

	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Str    string
	Number int
}

type HandlerDB struct {
	DB     *sql.DB
	Str    string
	Number int
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

func (h *HandlerDB) Index(w http.ResponseWriter, r *http.Request) {

	if !PathMethod("/", "GET", w, r) {
		return
	}

	rows, err := h.DB.Query("SELECT * FROM posts INNER JOIN users ON users.id = posts.user_id ORDER BY posts.id DESC;")
	if err != nil {
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	var Post []Posts
	var id_users, id_posts, user_login, user_password, user_email, post_date, post_user_id, post_category_id, post string
	for rows.Next() {
		err = rows.Scan(&id_posts, &post_date, &post_user_id, &post_category_id, &post, &id_users, &user_login, &user_password, &user_email)
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

func (h *HandlerDB) Registration(w http.ResponseWriter, r *http.Request) {
	fmt.Println("asdf")
	if !PathMethod("/registration/", "GET", w, r) {
		return
	}
	ExecTemp("templates/registration.html", "registration.html", &Error{Str: h.Str, Number: h.Number}, w, r)
	h.Str, h.Number = "", 0
}

func (h *HandlerDB) SignIn(w http.ResponseWriter, r *http.Request) {
	if !PathMethod("/auth/", "GET", w, r) {
		return
	}

	login := r.FormValue("login")
	if login == "" {
		ExecTemp("templates/signin.html", "signin.html", Error{Str: "Enter login", Number: 403}, w, r)
		return
	}

	row := h.DB.QueryRow("SELECT login, password FROM users WHERE users.login= ?;", login)
	var tempPassword []byte
	var tempLogin string
	err1 := row.Scan(&tempLogin, &tempPassword)
	if err1 != nil && err1 == sql.ErrNoRows {
		ExecTemp("templates/signin.html", "signin.html", Error{"Incorrect login", 403}, w, r)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		ExecTemp("templates/signin.html", "signin.html", Error{"Enter password", 403}, w, r)
		return
	}

	err := bcrypt.CompareHashAndPassword(tempPassword, []byte(password))
	if err != nil {
		ExecTemp("templates/signin.html", "signin.html", Error{"Invalid password", 403}, w, r)
		return
	}

	sessionID := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionID.String(),
	}
	http.SetCookie(w, cookie)

	ExecTemp("templates/index.html", "index.html", h.DB, w, r)
}

func (h *HandlerDB) Created(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if !PathMethod("/registration/created/", "POST", w, r) {
		return
	}

	login := r.FormValue("login")
	if login == "" {
		// h.Registration(w, r)
		ExecTemp("templates/registration.html", "registration.html", Error{"Enter login", 400}, w, r)
		return
	}

	if !test("login", login, "This login already exists", 400, h, w, r) {
		return
	}

	password := r.FormValue("password")
	if password == "" {
		ExecTemp("templates/registration.html", "registration.html", Error{"Enter password", 400}, w, r)
		return
	}

	email := r.FormValue("email")
	if !isEmailValid(email) {
		ExecTemp("templates/registration.html", "registration.html", Error{"Invalid email (everyone@example.com)", 400}, w, r)
		return
	}

	if !test("email", email, "This email already exists", 400, h, w, r) {
		return
	}

	GenPassword := []byte(password)
	if !GeneratePass(&GenPassword, w, r) {
		return
	}

	_, err := h.DB.Exec("INSERT INTO users (login, password, email) VALUES ( ?, ?, ?);", login, GenPassword, email)
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

func test(NameColumn, ValueColumn, str string, number int, h *HandlerDB, w http.ResponseWriter, r *http.Request) bool {
	row := h.DB.QueryRow("SELECT ? FROM users WHERE users."+NameColumn+"= ?;", NameColumn, ValueColumn)

	err := row.Scan()
	fmt.Println(err)
	fmt.Println(sql.ErrNoRows)
	if !errors.Is(err, sql.ErrNoRows) {
		fmt.Println("heloo")
		fmt.Println(str)
		h.Str, h.Number = str, number
		r.URL.Path = "/registration/"
		http.Redirect(w, r, r.Header.Get("/registration/"), 302)
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

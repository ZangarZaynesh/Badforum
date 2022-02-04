package structure

import (
	"database/sql"
	"errors"
	"fmt"
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

	// row := h.DB.QueryRow("SELECT ? FROM sessions WHERE ."+NameColumn+"= ?;", NameColumn, ValueColumn)

	Session := r.Cookies()
	fmt.Println("cookie: ", Session)
	// if err != nil {
	// 	cookie := &http.Cookie{
	// 		Name:    "session",
	// 		Value:   "failTest",
	// 		Expires: time.Now().Add(720 * time.Second),
	// 	}

	// 	http.SetCookie(w, cookie)
	// 	fmt.Println("WoooooW")
	// 	fmt.Println(Session)
	// 	// fmt.Println(Session.Name)
	// 	// fmt.Println(Session.Value)
	// } else {
	// 	fmt.Println(Session.Name)
	// 	fmt.Println(Session.Value)
	// 	fmt.Println(err)
	// }

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
	if !PathMethod("/registration/", "GET", w, r) {
		return
	}
	ExecTemp("templates/registration.html", "registration.html", &Error{Str: h.Str, Number: h.Number}, w, r)
	h.Str, h.Number = "", 0
}

func (h *HandlerDB) Created(w http.ResponseWriter, r *http.Request) {
	if !PathMethod("/registration/created/", "POST", w, r) {
		return
	}

	login := r.FormValue("login")
	if login == "" {
		h.Str, h.Number = "Enter login", 400
		http.Redirect(w, r, "/registration/", 302)
		return
	}

	if !test("login", login, "This login already exists", 400, h, w, r) {
		return
	}

	password := r.FormValue("password")
	if password == "" {
		h.Str, h.Number = "Enter password", 400
		http.Redirect(w, r, "/registration/", 302)
		return
	}

	email := r.FormValue("email")
	if !isEmailValid(email) {
		h.Str, h.Number = "Incorrected email", 400
		http.Redirect(w, r, "/registration/", 302)
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

func (h *HandlerDB) SignIn(w http.ResponseWriter, r *http.Request) {
	if !PathMethod("/auth/", "GET", w, r) {
		return
	}

	ExecTemp("templates/signIn.html", "signIn.html", h, w, r)
	h.Str, h.Number = "", 0
}

func (h *HandlerDB) SignAccess(w http.ResponseWriter, r *http.Request) {
	if !PathMethod("/auth/user/", "POST", w, r) {
		return
	}

	login := r.FormValue("login")
	if login == "" {
		h.Str, h.Number = "Enter login", 400
		http.Redirect(w, r, "/auth/", 302)
		return
	}

	row := h.DB.QueryRow("SELECT id, login, password FROM users WHERE users.login= ?;", login)
	var user_id int
	var tempPassword []byte
	var tempLogin string
	err1 := row.Scan(&user_id, &tempLogin, &tempPassword)
	if errors.Is(err1, sql.ErrNoRows) {
		h.Str, h.Number = "Incorrected login or password", 400
		http.Redirect(w, r, "/auth/", 302)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		h.Str, h.Number = "Enter password", 400
		http.Redirect(w, r, "/auth/", 302)
		return
	}
	BytePass := []byte(password)
	err := bcrypt.CompareHashAndPassword(tempPassword, BytePass)
	if err != nil {
		h.Str, h.Number = "Invalid password", 400
		http.Redirect(w, r, "/auth/", 302)
		return
	}

	sessionID := uuid.NewV4()
	cookie := &http.Cookie{
		Name:   "session",
		Value:  sessionID.String(),
		MaxAge: 300,
	}
	http.SetCookie(w, cookie)

	_, err = h.DB.Exec("INSERT INTO sessions (user_id, key) VALUES ( ?, ?);", user_id, sessionID)

	if err != nil {
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}
	// ExecTemp("templates/index.html", "index.html", h, w, r)
	// http.Redirect(w, r, "/", 302)

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
	if !errors.Is(err, sql.ErrNoRows) {
		h.Str, h.Number = str, number
		// r.URL.Path = "/registration/"
		http.Redirect(w, r, "/registration/", 302)
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

	// tmpl, err := template.ParseFiles(PathHTML)
	// if err != nil {
	// 	Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
	// 	return
	// }

	err := tmpl.ExecuteTemplate(w, NameHTML, Struct)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

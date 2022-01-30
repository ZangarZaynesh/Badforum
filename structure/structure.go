package structure

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
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

	if !PathMethod("/", "GET", w, r) {
		return
	}

	var Post []Posts

	rows, err := h.Art.Query("SELECT * FROM users INNER JOIN posts ON users.id = posts.user_id")
	if err != nil {
		log.Fatal(err)
	}
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

	fmt.Println("Begin")
	login := r.FormValue("login")
	if login == "" {
		fmt.Println("login==")
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if !test("login", login, h, w, r) {
		return
	}
	fmt.Println("loginTest")
	password := r.FormValue("password")

	if password == "" {
		fmt.Println("password==")
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	email := r.FormValue("email")

	if !isEmailValid(email) {
		fmt.Println("emailValid")
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
	}

	if !test("email", email, h, w, r) {
		return
	}
	fmt.Println("EmailTest")

	_, err := h.Art.Exec("INSERT INTO users (login, password, email) VALUES ( ?, ?, ?);", login, password, email)

	if err != nil {
		fmt.Println("ExecInsert")
		Err("400 Bad Request", http.StatusBadRequest, w, r)
		return
		// log.Fatalf("This error is in Created().InsertIntoUsers!!! %v", err)
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
	rows, err1 := h.Art.Query("SELECT ? FROM users WHERE users."+NameColumn+"= ?;", NameColumn, ValueColumn)

	if err1 != nil {
		fmt.Println("testQueryErr", err1)
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return false
	}

	var temp string
	for rows.Next() {
		err1 = rows.Scan(&temp)

		if err1 != nil {
			fmt.Println("testScanErr1")
			log.Fatal(err1)
		}

		if temp == ValueColumn {
			fmt.Println("test==ValueColumn")
			Err("400 Bad Request", http.StatusBadRequest, w, r)
			return false
		}
	}
	defer rows.Close()
	return true
}

func ExecTemp(PathHTML, NameHTML string, Struct interface{}, w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(PathHTML)
	if err != nil {
		fmt.Println("errExecTempTemplateParse")
		Err("500 Internal Server Error", http.StatusInternalServerError, w, r)
		return
	}

	err = tmpl.ExecuteTemplate(w, NameHTML, Struct)
	if err != nil {
		fmt.Println("errExecTempExecuteTemplate")
		log.Println("Error when parsing a template:", err)
		fmt.Fprintf(w, err.Error())
		return
	}
}

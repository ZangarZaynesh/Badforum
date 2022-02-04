package server

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"student/structure"
)

func NewHandlerDB(p *sql.DB) *structure.HandlerDB {
	return &structure.HandlerDB{DB: p}
}

func Routers(db *sql.DB) {
	tmpl, _ := template.ParseGlob("templates/*.html")
	handler := NewHandlerDB(db)
	// cookie, err := r.Cookie("session")
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/registration/", handler.Registration)
	http.HandleFunc("/registration/created/", handler.Created)
	http.HandleFunc("/auth/", handler.SignIn)
	http.HandleFunc("/auth/user/", handler.SignAccess)
	log.Println(http.ListenAndServe(":8080", nil))
}

package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"student/structure"
)

func NewHandlerDB(p *sql.DB) *structure.HandlerDB {
	return &structure.HandlerDB{DB: p}
}

func Routers(db *sql.DB) {
	handler := NewHandlerDB(db)
	fmt.Println(http.MethodGet)
	fmt.Println(http.MethodHead)
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/registration/", handler.Registration)
	http.HandleFunc("/registration/created/", handler.Created)
	http.HandleFunc("/auth/", handler.SignIn)
	log.Println(http.ListenAndServe(":8080", nil))
}

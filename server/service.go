package server

import (
	"database/sql"
	"log"
	"net/http"
	"student/structure"
)

func NewHandlerArt(p *sql.DB) *structure.HandlerArt {
	return &structure.HandlerArt{p}
}

func Routers(db *sql.DB) {
	handler := NewHandlerArt(db)
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/registration", handler.Registration)
	http.HandleFunc("/registration/created", handler.Created)
	// http.HandleFunc("/filters/", Filter)
	log.Println(http.ListenAndServe(":8080", nil))
}

package main

import (
	"fmt"
	"net/http"
	"student/db"
	"student/server"
)

func main() {
	db := db.CheckDB()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("http://localhost:8080/ is listening....")
	server.Routers(db)
	defer db.Close()
}

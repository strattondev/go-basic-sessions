package main

import (
	"github.com/strattonw/go-basic-sessions/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/authenticated", handlers.Authenticated)
	http.HandleFunc("/logout", handlers.Logout)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

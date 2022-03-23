package main

import (
	"log"
	"net/http"

	"github.com/yimikao/googleoauth/handlers"
)

func main() {
	sv := http.Server{
		Addr:    ":8080",
		Handler: handlers.NewMux(),
	}

	log.Printf("Starting HTTP Server. Listening at %q", sv.Addr)
	if err := sv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}

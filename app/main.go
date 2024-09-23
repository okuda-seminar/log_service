package main

import (
	"log"
	"log_service/app/infrastructure/mysql/db"
	"net/http"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatalf("db.Connect: %v", err)
	}
	log.Println("db connected")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			log.Fatalf("db.Ping: %v", err)
		}
	})
	log.Fatalf("http.ListenAndServe: %v", http.ListenAndServe(":8080", nil))
}

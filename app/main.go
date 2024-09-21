package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World\n")); err != nil {
			log.Fatalf("w.Write: %v", err)
		}
	})
	log.Fatalf("http.ListenAndServe: %v", http.ListenAndServe(":8080", nil))
}

package main

import (
	"log"
	"log_service/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

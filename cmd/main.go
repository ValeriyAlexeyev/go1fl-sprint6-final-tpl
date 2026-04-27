package main

import (
	"log"
	"os"

	"github.com/ValeriyAlexeyev/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	srv := server.NewServer(logger)

	logger.Println("server started on :8080")

	if err := srv.HTTP.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

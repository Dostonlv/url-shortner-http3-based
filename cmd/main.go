package main

import (
	"log"
	"os"

	"github.com/Dostonlv/url-shortner-http3-based/internal/handler"
	"github.com/Dostonlv/url-shortner-http3-based/internal/repository/postgres"
	"github.com/Dostonlv/url-shortner-http3-based/internal/service"
	"github.com/Dostonlv/url-shortner-http3-based/pkg/database"
	"github.com/joho/godotenv"
	"github.com/quic-go/quic-go/http3"
)

// @title URL Shortener API
// @version 1.0
// @description HTTP/3 based URL shortener service with PostgreSQL
// @host localhost:4433
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db, err := database.NewPostgresDB(database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories, services and handlers
	repos := postgres.NewURLRepository(db)
	services := service.NewURLService(repos)
	handlers := handler.NewHandler(services)

	// Start server
	log.Println("Server running on https://localhost:4433")
	err = http3.ListenAndServeTLS(
		":4433",
		"certs/server.crt",
		"certs/server.key",
		handlers.InitRoutes(),
	)
	if err != nil {
		log.Fatal(err)
	}
}

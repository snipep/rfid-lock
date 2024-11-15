package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/snipep/iot/internal/database"
	"github.com/snipep/iot/internal/router"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	// Load environment variables
	err := godotenv.Load("./../../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Parse port
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Invalid PORT. Defaulting to 8080")
		// port = 8080
	}

	// Initialize database
	database.InitDB()

	// Create and return HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.RegisterRoutes(),
	}

	return server
}

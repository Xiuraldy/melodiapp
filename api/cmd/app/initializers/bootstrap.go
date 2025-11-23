package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"melodiapp/cmd/app/routes"
)

func Run() {
	if os.Getenv("GIN_MODE") != "release" {
		if err := godotenv.Load(); err != nil {
			log.Printf("No .env file loaded (godotenv): %v. Continuing with existing environment variables.", err)
		}
	}

	InitDatabase()

	router := routes.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

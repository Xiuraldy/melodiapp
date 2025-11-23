package initializers

import (
	"log"

	"melodiapp/database"
	"melodiapp/models"
)

func InitDatabase() {
	log.Println("Initializing database connection...")
	database.CreateDbConnection()

	if err := database.DBConn.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}
	log.Println("Database initialized and migrations applied")
}

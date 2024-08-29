package main

import (
	"goproject/config"
	"goproject/database"
	"goproject/models"
	"goproject/routes"
	"log"
	"os"
)

func main() {
	env := os.Getenv("GIN_ENV")
	if env == "" {
		env = "dev"
	}
	config.LoadConfig(env)

	database.ConnectDatabase()

	// Run database migrations
	err := database.DB.AutoMigrate(&models.GoUser{}, &models.GoArticle{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := routes.SetupRouter()
	r.Run(":8000")
}

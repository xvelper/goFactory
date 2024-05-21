package main

import (
	"goFactory/controllers"
	"goFactory/handlers"
	"goFactory/models"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize database
	db := models.InitDB()

	// Create repos directory if not exists
	if _, err := os.Stat("repos"); os.IsNotExist(err) {
		os.Mkdir("repos", 0755)
	}

	// Initialize Fiber
	app := fiber.New()

	// Setup routes
	app.Delete("/repos/:username/:repo", handlers.DeleteRepositoryHandler)
	app.Post("/users", handlers.CreateUserHandler(db))
	app.Post("/repos/init", handlers.InitRepositoryHandler)
	app.Get("/repos", handlers.GetRepoHandler(db))

	// Basic authentication middleware with database
	auth := controllers.BasicAuthMiddleware(db)

	// Git HTTP backend with authentication
	app.All("/repos/:username/:repo/info/refs", auth, handlers.GitHTTPHandler)

	// Start server
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}

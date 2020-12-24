package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"implementingChange/requestHandlers"
	"log"
)

// assignRoutes maps web server routes to the handler functions.
func assignRoutes(app *fiber.App) {
	app.Get("/", requestHandlers.IndexHandler)
	app.Post("/login", requestHandlers.LoginHandler)
}

// loadEnv loads environment variables from `.env`.
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	app := fiber.New()
	assignRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

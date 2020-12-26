package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"implementingChange/auth"
	"implementingChange/db"
	"implementingChange/handlers"
	"log"
)

func injectMiddleware(app *fiber.App) {
	app.Use("/v1", auth.CheckAuth)
}

// assignRoutes maps web server routes to the handler functions.
func assignRoutes(app *fiber.App) {
	app.Post("/login", handlers.LoginHandler)
	app.Get("/v1/ping", handlers.PingHandler)
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
	if err := db.InitRedis(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	injectMiddleware(app)
	assignRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

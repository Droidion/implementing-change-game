package main

import (
	"github.com/Droidion/implementing-change-game/db"
	"github.com/Droidion/implementing-change-game/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
)

const apiVer = "/v1"

// injectMiddleware assigns middleware handlers to certain paths
func injectMiddleware(app *fiber.App) {
	app.Use(apiVer, handlers.CheckAuthMiddleware)
}

// assignRoutes maps web server routes to the handler functions
func assignRoutes(app *fiber.App) {
	app.Post("/login", handlers.LoginHandler)
	app.Get(apiVer+"/ping", handlers.PingHandler)
	app.Post(apiVer+"/logout", handlers.LogoutHandler)
}

// loadEnv loads environment variables from .env file
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
	db.PostgresConnect()
	defer db.PgConn.Close()
	app := fiber.New()
	injectMiddleware(app)
	assignRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

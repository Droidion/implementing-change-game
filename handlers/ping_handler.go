package handlers

import "github.com/gofiber/fiber/v2"

// PingHandler handles GET request for dummy page
func PingHandler(c *fiber.Ctx) error {
	return c.SendString("Pong")
}

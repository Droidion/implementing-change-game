package requestHandlers

import "github.com/gofiber/fiber/v2"

// IndexHandler handles GET request for root page
func IndexHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

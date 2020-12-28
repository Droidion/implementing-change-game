package handlers

import (
	"github.com/Droidion/implementing-change-game/db"
	"github.com/Droidion/implementing-change-game/models"
	"github.com/gofiber/fiber/v2"
)

// LogoutHandler handles POST requests for /logout
func LogoutHandler(c *fiber.Ctx) error {
	// Get access token metadata from fasthttp context
	// It's supposed to be injected by middleware that checks authorization
	accessDetails, ok := c.Context().UserValue("access_details").(*models.AccessDetails)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "Something wrong with auth details in token")
	}

	// Delete access token metadata from redis
	if _, err := db.DeleteAuth(accessDetails.AccessUuid); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Data provided in token was not found on server, token probably expired")
	}

	return c.JSON("Successfully logged out")
}

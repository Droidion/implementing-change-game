package handlers

import (
	"github.com/Droidion/implementing-change-game/auth"
	"github.com/Droidion/implementing-change-game/db"
	"github.com/gofiber/fiber/v2"
)

// CheckAuthMiddleware is a Fiber middleware that checks if the token in header
func CheckAuthMiddleware(c *fiber.Ctx) error {
	// Extract token from header
	tokenStr, err := auth.ExtractToken(c.Request().Header.Peek("Authorization"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "No Authorization header provided")
	}

	// Verify token content
	accessDetails, err := auth.VerifyToken(tokenStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Token not verified")
	}

	// Get token metadata cached in Redis
	_, err = db.FetchAuth(accessDetails)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Token is correct, but the server does not have it no more, it's probably expired")
	}

	// Save access token metadata in fasthttp request context for reusing in the following request handlers
	c.Context().SetUserValue("access_details", accessDetails)

	// Move to the next request handler
	return c.Next()
}

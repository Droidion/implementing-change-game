package handlers

import (
	"github.com/Droidion/implementing-change-game/auth"
	"github.com/Droidion/implementing-change-game/models"
	"github.com/gofiber/fiber/v2"
)

// loginRequest contains request body for POST /login
type loginRequest struct {
	Password string `json:"password" xml:"password" form:"password"`
}

// loginResponse contains response body for POST /login
type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LoginHandler handles POST requests for /login
func LoginHandler(c *fiber.Ctx) error {
	var err error

	// Dummy user
	var user = models.User{
		Id:       123,
		Username: "username",
		Password: "password",
	}

	// Parse request body
	requestBody := new(loginRequest)
	if err := c.BodyParser(requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Incorrect request body")
	}

	// Check that password is correct (with dummy user)
	if user.Password != requestBody.Password {
		return fiber.NewError(fiber.StatusUnauthorized, "Could not recognize the requestBody")
	}

	// Create token
	tokenDetails, err := auth.CreateToken(user.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not create tokens")
	}

	// Cache token
	if err = auth.CacheTokens(user.Id, tokenDetails); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not cache tokens")
	}

	// Make response
	responseBody := loginResponse{tokenDetails.AccessToken, tokenDetails.RefreshToken}
	return c.JSON(responseBody)
}

package handlers

import (
	"github.com/Droidion/implementing-change-game/auth"
	"github.com/Droidion/implementing-change-game/db"
	"github.com/Droidion/implementing-change-game/models"
	"github.com/gofiber/fiber/v2"
)

// loginRequest contains request body for POST /login
type loginRequest struct {
	Login    string `json:"login" xml:"login" form:"login"`
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

	// Parse request body
	requestBody := new(loginRequest)
	if err = c.BodyParser(requestBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Incorrect request body")
	}

	var user *models.User
	user, err = db.GetUserByLogin(requestBody.Login)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Login not found")
	}

	_, err = auth.CompareHashAndPassword(user.Password, requestBody.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "IncorrectPassword")
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

package requestHandlers

import (
	"github.com/gofiber/fiber/v2"
	"implementingChange/auth"
	"implementingChange/models"
)

// LoginRequest contains request body for POST /login
type loginRequest struct {
	Password string `json:"password" xml:"password" form:"password"`
}

// LoginResponse contains response body for POST /login
type loginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles POST requests for /login
func LoginHandler(c *fiber.Ctx) error {
	var user = models.User{
		ID:       1,
		Username: "username",
		Password: "password",
	}

	password := new(loginRequest)
	if err := c.BodyParser(password); err != nil {
		return err
	}

	if user.Password == password.Password {
		token, err := auth.CreateToken(123)
		if err != nil {
			return err
		}

		return c.JSON(loginResponse{token})
	}

	return fiber.NewError(fiber.StatusUnauthorized, "Could not recognize the password")
}
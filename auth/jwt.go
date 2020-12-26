// Package auth implements methods for using in authentication and authorization pipelines

package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"implementingChange/db"
	"implementingChange/models"
	"os"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

// createUuid creates new random UUID (v4)
func createUuid() (string, error) {
	if id, err := uuid.NewRandom(); err != nil {
		return "", eris.Wrap(err, "could not create UUID")
	} else {
		return id.String(), nil
	}
}

// generateTokenFromClaims generates and signs the token from a set of claims
func generateTokenFromClaims(claims jwt.MapClaims, secret string) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte(secret))
	if err != nil {
		return "", eris.Wrap(err, "could not sign the token")
	}
	return signedToken, nil
}

// CreateToken creates a struct with token metadata
func CreateToken(userid uint64) (*models.TokenDetails, error) {
	var err error

	tokenDetails := &models.TokenDetails{}

	// Add expiration times and UUIDs
	// Access token should expire in 15 minutes
	tokenDetails.AccessExpires = time.Now().Add(time.Minute * 15).Unix()
	if tokenDetails.AccessUuid, err = createUuid(); err != nil {
		return nil, eris.Wrap(err, "could not create Access UUID")
	}
	// Refresh token should expire in 7 days
	tokenDetails.RefreshExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	if tokenDetails.RefreshUuid, err = createUuid(); err != nil {
		return nil, eris.Wrap(err, "could not create Refresh UUID")
	}

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = tokenDetails.AccessUuid
	accessTokenClaims["user_id"] = userid
	accessTokenClaims["exp"] = tokenDetails.AccessExpires
	tokenDetails.AccessToken, err = generateTokenFromClaims(accessTokenClaims, os.Getenv("ACCESS_SECRET"))
	if err != nil {
		return nil, eris.Wrap(err, "could not create Access Token")
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userid
	refreshTokenClaims["exp"] = tokenDetails.RefreshExpires
	tokenDetails.RefreshToken, err = generateTokenFromClaims(refreshTokenClaims, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		return nil, eris.Wrap(err, "could not create Refresh Token")
	}

	return tokenDetails, nil
}

// saveTokenToRedis saves a single token to Redis
func saveTokenToRedis(expires int64, uuid string, userId uint64) error {
	utcTime := time.Unix(expires, 0)
	now := time.Now()
	err := db.Redis.Set(ctx, "impl-change-token-" + uuid, strconv.Itoa(int(userId)), utcTime.Sub(now)).Err()
	if err != nil {
		return eris.Wrap(err, "could not save token to Redis")
	}
	return nil
}

// CacheTokens saves tokens to Redis
func CacheTokens(userId uint64, tokenDetails *models.TokenDetails) error {
	var err error
	if err = saveTokenToRedis(tokenDetails.AccessExpires, tokenDetails.AccessUuid, userId); err != nil {
		return eris.Wrap(err, "could not cache access token")
	}
	if err = saveTokenToRedis(tokenDetails.RefreshExpires, tokenDetails.RefreshUuid, userId); err != nil {
		return eris.Wrap(err, "could not cache refresh token")
	}
	return nil
}

func extractToken(authHeaderVal []byte) (string, error) {
	authHeaderValArr := strings.Split(string(authHeaderVal), " ")
	if len(authHeaderValArr) != 2 {
		return "", eris.New("Could not parse Authorization header")
	}
	return authHeaderValArr[1], nil
}

func verifyToken(tokenStr string) (*models.AccessDetails, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, eris.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, eris.New("could not parse the token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, eris.New("no access_uuid in token")
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, eris.New("incorrect user_id in token")
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:   userId,
		}, nil
	}
	return nil, eris.New("incorrect token content")
}

func fetchAuth(authD *models.AccessDetails) (uint64, error) {
	userid, err := db.Redis.Get(ctx, "impl-change-token-" + authD.AccessUuid).Result()
	if err != nil {
		return 0, eris.Wrap(err, "could not find token in Redis")
	}
	userID, err := strconv.ParseUint(userid, 10, 64)
	if err != nil {
		return 0, eris.Wrap(err, "could not parse user id")
	}
	return userID, nil
}

// CheckAuth is a Fiber middleware that checks if the token in header
func CheckAuth(c *fiber.Ctx) error {
	tokenStr, err := extractToken(c.Request().Header.Peek("Authorization"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "No Authorization header provided")
	}

	accessDetails, err := verifyToken(tokenStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Token not verified")
	}

	_, err = fetchAuth(accessDetails)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Token is correct, but the server does not have it no more, it's probably expired")
	}

	return c.Next()
}
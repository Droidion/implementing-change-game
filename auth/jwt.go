// Package auth implements methods for using in authentication and authorization pipelines

package auth

import (
	"fmt"
	"github.com/Droidion/implementing-change-game/db"
	"github.com/Droidion/implementing-change-game/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"os"
	"strconv"
	"strings"
	"time"
)

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

	// Generate access token with claims
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = tokenDetails.AccessUuid
	accessTokenClaims["user_id"] = userid
	accessTokenClaims["exp"] = tokenDetails.AccessExpires
	tokenDetails.AccessToken, err = generateTokenFromClaims(accessTokenClaims, os.Getenv("ACCESS_SECRET"))
	if err != nil {
		return nil, eris.Wrap(err, "could not create Access Token")
	}

	// Generate refresh token with claims
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

// CacheTokens saves tokens to Redis
func CacheTokens(userId uint64, tokenDetails *models.TokenDetails) error {
	var err error

	// Save access token
	if err = db.SaveTokenToRedis(tokenDetails.AccessExpires, tokenDetails.AccessUuid, userId); err != nil {
		return eris.Wrap(err, "could not cache access token")
	}

	// Save refresh token
	if err = db.SaveTokenToRedis(tokenDetails.RefreshExpires, tokenDetails.RefreshUuid, userId); err != nil {
		return eris.Wrap(err, "could not cache refresh token")
	}

	return nil
}

// ExtractToken extracts token substring from Authorization header
func ExtractToken(authHeaderVal []byte) (string, error) {
	authHeaderValArr := strings.Split(string(authHeaderVal), " ")
	if len(authHeaderValArr) != 2 {
		return "", eris.New("Could not parse Authorization header")
	}
	return authHeaderValArr[1], nil
}

// VerifyToken checks if token is legit and extracts metadata from it
func VerifyToken(tokenStr string) (*models.AccessDetails, error) {
	// Parse the token and check that it has proper signing method
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, eris.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, eris.New("could not parse the token")
	}

	// Checks the token metadata
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
			UserId:     userId,
		}, nil
	}

	return nil, eris.New("incorrect token content")
}

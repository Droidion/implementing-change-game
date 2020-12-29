package auth

import (
	"github.com/andskur/argon2-hashing"
	"github.com/rotisserie/eris"
)

// EncryptPassword encrypts plain password using argon2id
func EncryptPassword(plainPassword string) (string, error) {
	hash, err := argon2.GenerateFromPassword([]byte(plainPassword), argon2.DefaultParams)
	if err != nil {
		return "", eris.Wrap(err, "could not encrypt password")
	}
	return string(hash), nil
}

// CompareHashAndPassword checks if hashed password is the same as plain password
func CompareHashAndPassword(hash string, plain string) (bool, error) {
	err := argon2.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false, eris.Wrap(err, "passwords do not match")
	}
	return true, nil
}

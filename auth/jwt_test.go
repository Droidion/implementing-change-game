package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateToken(t *testing.T) {
	token, err := CreateToken(123)
	assert.Greater(t, len(token), 0, "token should be non empty string")
	assert.Equal(t, err, nil, "there should be no error while generating token")
}

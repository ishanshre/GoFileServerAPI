package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generatePassword(t *testing.T) (string, string) {
	plainPassword := RandomAlphapet(10)
	hashed, err := GeneratePassword(plainPassword)
	assert.NoError(t, err)
	assert.NotEqual(t, hashed, "")
	return hashed, plainPassword
}

func TestGeneratePassword(t *testing.T) {
	generatePassword(t)
}

func TestComparePassword(t *testing.T) {
	hashed, plain := generatePassword(t)
	err := ComparePassword(hashed, plain)
	assert.NoError(t, err)
	err = ComparePassword(hashed, "asdasd")
	assert.Error(t, err)
}

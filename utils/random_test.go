package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	min := int64(0)
	max := int64(100)

	result := RandomInt(min, max)
	assert.GreaterOrEqual(t, result, min)
	assert.LessOrEqual(t, result, max)
}

func TestRandomString(t *testing.T) {
	length := 10
	resutl := RandomAlphapet(10)
	assert.Equal(t, length, len(resutl))
}
func TestRandomNumString(t *testing.T) {
	length := 10
	resutl := RandomAlphaNum(10)
	assert.Equal(t, length, len(resutl))
}

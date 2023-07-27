package utils

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const alphaNum = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomAlphapet(length int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomAlphaNum(length int) string {
	var sb strings.Builder
	k := len(alphaNum)
	for i := 0; i < length; i++ {
		c := alphaNum[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

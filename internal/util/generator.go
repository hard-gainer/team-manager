package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecureToken(lenght int) string {
	b := make([]byte, lenght)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

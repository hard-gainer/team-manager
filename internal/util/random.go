package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	chars   = "abcdefghijklmnopqrstuvwxyz"
	nums    = "0123456789"
	alpabet = chars + nums
)

func init() {
	rand.New(rand.NewSource(time.Now().Unix()))
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alpabet)

	for i := 0; i < n; i++ {
		c := alpabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String() 
}

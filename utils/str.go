package utils

import (
	"math/rand"
	"strings"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func RandLowStr(n int) string {
	return strings.ToLower(RandStr(n))
}

func RandUpperStr(n int) string {
	return strings.ToUpper(RandStr(n))
}

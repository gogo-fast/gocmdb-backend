package utils

import (
	"math/rand"
	"strings"
	"time"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
const cCount = len(characters)

func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	var chars = make([]string, length)
    cs := characters[:]
	for i := 0; i < length; i++ {
		chars[i] = string(cs[rand.Intn(cCount)])
	}
	return strings.Join(chars,"")
}

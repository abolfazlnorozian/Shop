package helpers

import (
	"math/rand"
	"time"
)

func GenerateRandomCode(length int) string {
	table := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = table[rand.Intn(len(table))]
	}
	return string(code)
}

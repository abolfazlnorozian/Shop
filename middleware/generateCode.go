package middleware

import (
	"math/rand"
	"time"
)

// func EncodeToString(max int) string {
// 	b := make([]byte, max)
// 	n, err := io.ReadAtLeast(rand.Reader, b, max)
// 	if n != max {
// 		panic(err)
// 	}
// 	for i := 0; i < len(b); i++ {
// 		b[i] = table[int(b[i])%len(table)]
// 	}
// 	return string(b)
// }

// var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateRandomCode(length int) string {
	table := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = table[rand.Intn(len(table))]
	}
	return string(code)
}

package helpers

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomUsername(phoneNumber string) string {
	// Remove non-digit characters from phoneNumber
	phoneNumber = removeNonDigitCharacters(phoneNumber)

	// Generate a random suffix for the username
	rand.Seed(time.Now().UnixNano())
	randomSuffix := rand.Intn(9999999999)

	// Combine the phoneNumber and random suffix to create the username
	username := "users" + strconv.Itoa(randomSuffix)

	return username
}

// Function to remove non-digit characters from a string
func removeNonDigitCharacters(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

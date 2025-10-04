package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano() + rand.Int63n(1000000))
}

// RandSeq generates a random string to serve as dummy data
func RandNumber(n int) string {

	letters := []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

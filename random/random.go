package random
import (
	"math/rand"
	"time"
)

// RandSeq generates a random string to serve as dummy data
func RandNumber(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func RandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func RandomString1(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

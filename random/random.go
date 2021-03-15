// Package random provides methods to generate random data
package random

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRandomInt generates a random integer min and max
func GetRandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// GetRandomString generates a random string of length n
func GetRandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// GetRandomStringListOfSizeN return a list of strings
func GetRandomStringListOfSizeN(n int) *[]string {
	var result []string
	for i := 0; i < n; i++ {
		temp := GetRandomString(6)
		result = append(result, temp)
	}
	return &result
}

// GetRandomUUID returns a random uuid string
func GetRandomUUID() string {
	return uuid.New().String()
}

package isunippets

import (
	"math/rand"
	"regexp"
)
import "github.com/google/uuid"

var regexpFormat = regexp.MustCompile("^is_dirty=(true|false)$")

func RegexpIsMatch(condition string) bool {
	return regexpFormat.MatchString(condition)
}

func GenerateUUID() string {
	return uuid.New().String()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

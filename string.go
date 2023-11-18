package isunippets

import "regexp"
import "github.com/google/uuid"

var regexpFormat = regexp.MustCompile("^is_dirty=(true|false)$")

func RegexpIsMatch(condition string) bool {
	return regexpFormat.MatchString(condition)
}

func GenerateUUID() string {
	return uuid.New().String()
}

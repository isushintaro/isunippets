package isunippets

import "regexp"

var regexpFormat = regexp.MustCompile("^is_dirty=(true|false)$")

func RegexpIsMatch(condition string) bool {
	return regexpFormat.MatchString(condition)
}

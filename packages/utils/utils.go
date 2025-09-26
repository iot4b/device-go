package utils

import "regexp"

func MatchRegex(pattern, str string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

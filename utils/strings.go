package utils

import (
	"regexp"
	"strings"
)

func Underscore(str string) string {
	str = strings.ReplaceAll(str, "::", "/")
	str = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`).ReplaceAllString(str, "${1}_${2}")
	str = regexp.MustCompile(`([a-z\d])([A-Z])`).ReplaceAllString(str, "${1}_${2}")
	str = strings.ReplaceAll(str, "-", "_")
	str = strings.ToLower(str)
	return str
}

func Plain(str string) string {
	str = strings.ReplaceAll(str, "\r\n", " ")
	str = regexp.MustCompile(`[\s\-_=]{2,}`).ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	return str
}

func CutN(str string, size int) string {
	runes := []rune(str)

	if len(runes) < size {
		return string(runes)
	}

	return string(runes[0:size])
}

func Cut(str string) string {
	return CutN(str, 50)
}

func PayloadType(str string) string {
	str = strings.ReplaceAll(str, "Event", "")
	str = regexp.MustCompile(`.*Comment`).ReplaceAllString(str, "Comment")
	str = strings.ReplaceAll(str, "Issues", "Issue")
	str = Underscore(str)
	return str
}

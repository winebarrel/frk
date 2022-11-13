package utils

import (
	"regexp"
	"strings"
)

type Str string

func (s Str) String() string {
	return string(s)
}

func (s Str) Underscore() Str {
	str := string(s)
	str = strings.ReplaceAll(str, "::", "/")
	str = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`).ReplaceAllString(str, "${1}_${2}")
	str = regexp.MustCompile(`([a-z\d])([A-Z])`).ReplaceAllString(str, "${1}_${2}")
	str = strings.ReplaceAll(str, "-", "_")
	str = strings.ToLower(str)
	return Str(str)
}

func (s Str) Plain() Str {
	str := string(s)
	str = strings.ReplaceAll(str, "\r\n", " ")
	str = regexp.MustCompile(`[\s\-_=]{2,}`).ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	return Str(str)
}

func (s Str) CutN(size int) Str {
	str := string(s)
	runes := []rune(str)

	if len(runes) < size {
		return Str(runes)
	}

	return Str(runes[0:size])
}

func (s Str) Cut() Str {
	return s.CutN(50)
}

func (s Str) PayloadType() Str {
	str := string(s)
	str = strings.ReplaceAll(str, "Event", "")
	str = regexp.MustCompile(`.*Comment`).ReplaceAllString(str, "Comment")
	str = strings.ReplaceAll(str, "Issues", "Issue")
	return Str(str).Underscore()
}

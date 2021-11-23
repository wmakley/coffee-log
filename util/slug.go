package util

import (
	"errors"
	"regexp"
	"strings"
)

var replacePattern *regexp.Regexp

func init() {
	replacePattern = regexp.MustCompile("[^a-zA-Z0-9\\-_]")
}

func Sluggify(input string) (string, error) {
	if input == "" {
		return input, errors.New("input must not be empty")
	}

	lower := strings.ToLower(input)
	return replacePattern.ReplaceAllString(lower, "-"), nil
}

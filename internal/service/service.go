package service

import (
	"errors"
	"strings"

	"github.com/ValeriyAlexeyev/go1fl-sprint6-final/pkg/morse"
)

var ErrEmptyInput = errors.New("input is empty")

func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", ErrEmptyInput
	}

	if isMorse(input) {
		return morse.ToText(input), nil
	}

	return morse.ToMorse(input), nil
}

func isMorse(s string) bool {
	for _, r := range s {
		if r != '.' &&
			r != '-' &&
			r != ' ' &&
			r != '/' &&
			r != '\n' &&
			r != '\r' &&
			r != '\t' {
			return false
		}
	}
	return true
}

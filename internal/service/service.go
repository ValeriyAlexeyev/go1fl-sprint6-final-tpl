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

	// если есть точка или тире → считаем что это Morse
	if strings.Contains(input, ".") || strings.Contains(input, "-") {
		return morse.ToText(input), nil
	}

	return morse.ToMorse(input), nil
}

package utils

import (
	"fmt"
	"strings"
)

func ParseToken(token string) (string, string, error) {
	parts := strings.SplitN(token, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid token: %s", token)
	}
	return strings.ToLower(parts[0]), parts[1], nil
}

func ConvertToBytes(size int, unit string) (int, error) {
	switch unit {
	case "B", "b":
		return size, nil
	case "K", "k":
		return size * 1024, nil
	case "M", "m":
		return size * 1024 * 1024, nil
	case "G", "g":
		return size * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("invalid unit: %s", unit)
	}
}

var alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y", "Z",
}

var pathToLetter = make(map[string]string)

var nextLetterIndex = 0

func GetLetter(path string) (string, error) {
	if _, exists := pathToLetter[path]; !exists {
		if nextLetterIndex < len(alphabet) {
			pathToLetter[path] = alphabet[nextLetterIndex]
			nextLetterIndex++
		} else {
			return "", fmt.Errorf("no more letters available")
		}
	}

	return pathToLetter[path], nil
}

package helpers

import (
	"errors"
	"strings"
)

func NormalizeTelephone(telephone string) (string, error) {
	telephone = strings.TrimSpace(telephone)
	if telephone == "" {
		return "", errors.New("telephone is required")
	}

	switch telephone[0] {
	case '0':
		telephone = "62" + telephone[1:]
	case '+':
		telephone = telephone[1:]
	case '8':
		telephone = "62" + telephone
	case '6':
	default:
		return "", errors.New("invalid telephone number format")
	}

	if len(telephone) < 2 {
		return "", errors.New("invalid telephone number format")
	}

	return telephone, nil
}

package user

import (
	"errors"
	"strings"
)

func validateName(s string) (string, error) {
	s = strings.TrimSpace(s)

	if len(strings.Split(s, " ")) != 2 {
		err := errors.New("должно быть 2 слова как на бейдже")

		return "", err
	}

	return s, nil
}

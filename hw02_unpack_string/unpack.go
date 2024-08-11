package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var (
		lastPos  = -1
		lastChar string
		sb       strings.Builder
	)

	for pos, char := range str {
		if unicode.IsDigit(char) {
			if lastPos < 0 {
				return "", ErrInvalidString
			}

			count, _ := strconv.Atoi(string(char))
			sb.WriteString(strings.Repeat(lastChar, count))

			lastPos = -1
			lastChar = ""
			continue
		}

		lastPos = pos
		if lastChar != "" {
			sb.WriteString(lastChar)
		}
		lastChar = string(char)
	}

	if lastChar != "" {
		sb.WriteString(lastChar)
	}

	return sb.String(), nil
}

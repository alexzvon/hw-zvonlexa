package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	var ss string
	sr := []rune(in)

	if err := validate(in); err != nil {
		return "", err
	}

	if in == "" {
		return ss, nil
	}

	for i := 0; i < len(sr); i++ {
		if i == len(sr)-1 {
			if !unicode.IsDigit(sr[i]) {
				ss += string(sr[i])
			}

			break
		}

		if sr[i+1] == '0' {
			i++
			continue
		}

		if unicode.IsDigit(sr[i+1]) {
			j, err := strconv.Atoi(string(sr[i+1]))
			if err != nil {
				return "", err
			}

			ss = concat(ss, strings.Repeat(string(sr[i]), j))

			i++
		} else {
			ss = concat(ss, string(sr[i]))
		}
	}

	return ss, nil
}

func validate(s string) error {
	var k int
	sr := []rune(s)

	if s == "" {
		return nil
	}

	if unicode.IsDigit(sr[0]) {
		return ErrInvalidString
	}

	k = 0
	for i := 0; i < len(sr); i++ {
		if unicode.IsDigit(sr[i]) {
			k++
		} else if k > 0 {
			k--
		}

		if k > 1 {
			return ErrInvalidString
		}
	}

	return nil
}

func concat(s ...string) string {
	var builder strings.Builder
	var count int

	for _, v := range s {
		count += len(v)
	}

	builder.Grow(count)

	for _, v := range s {
		builder.WriteString(v)
	}

	return builder.String()
}

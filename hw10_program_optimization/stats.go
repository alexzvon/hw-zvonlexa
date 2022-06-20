package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

const (
	dog = byte('@')
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, errors.New("domain must not be empty")
	}

	var user User

	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()

		err := user.UnmarshalJSON(line)
		if err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, domain) {
			n := strings.IndexByte(user.Email, dog) + 1

			result[strings.ToLower(user.Email[n:])]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("shouldn't see an error scanning a string")
	}

	return result, nil
}

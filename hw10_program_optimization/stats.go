package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, errors.New("domain must not be empty")
	}

	var user User

	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigFastest

	for scanner.Scan() {
		line := scanner.Bytes()
		err := json.Unmarshal(line, &user)
		if err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("shouldn't see an error scanning a string")
	}

	return result, nil
}

package hw10programoptimization

import (
	"bufio"
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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	suffix := "." + domain
	var user User
	for scanner.Scan() {
		user.Email = ""
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, suffix) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}

package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
	"strings"
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

func GetDomainStat(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	pusers := parseUsers(r)
	for pu := range pusers {
		if pu.err != nil {
			err = pu.err
			return
		}
		l1, l2 := getDomain(pu.value.Email)
		if l1 == domain {
			n := result[l2]
			n++
			result[l2] = n
		}
	}
	return
}

type parsedUser struct {
	value User
	err   error
}

func parseUsers(r io.Reader) (resultCh chan parsedUser) {
	resultCh = make(chan parsedUser, 10)

	go func() {
		defer close(resultCh)
		var uline []byte
		var user User
		var err error

		br := bufio.NewReader(r)
		for err == nil {
			uline, err = br.ReadBytes('}')
			if err != nil {
				break
			}
			err = json.Unmarshal(uline, &user)
			if err != nil {
				break
			}
			resultCh <- parsedUser{user, nil}
		}
		if err != io.EOF {
			resultCh <- parsedUser{User{}, err}
		}
	}()
	return
}

var regexpDomain *regexp.Regexp

func getDomain(email string) (level1, level2 string) {
	if regexpDomain == nil {
		regexpDomain = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+([\w-]{2,4})$`)
	}
	res := regexpDomain.FindStringSubmatch(strings.ToLower(email))
	if res == nil {
		return "", ""
	}
	return res[2], res[1] + res[2]
}

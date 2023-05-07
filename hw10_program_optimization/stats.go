package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	qjson "github.com/buger/jsonparser"
)

type User struct {
	//	ID       int
	//	Name     string
	//	Username string
	Email string
	// Phone    string
	// Password string
	// Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	err = parseRecords(r, &result, domain)
	return
}

func parseRecords(r io.Reader, domainStat *DomainStat, lvl1 string) error {
	var uline []byte
	var user User
	var err error

	br := bufio.NewReader(r)
	for err == nil {
		uline, err = br.ReadBytes('}')
		if err != nil {
			break
		}
		user.Email, err = qjson.GetString(uline, "Email")
		if err != nil {
			break
		}
		user.Email = strings.ToLower(user.Email)
		err = checkDomain(lvl1, user.Email, domainStat)
		if err != nil {
			break
		}
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func checkDomain(lvl1, email string, domainStat *DomainStat) error {
	arr1 := strings.Split(email, "@")
	if len(arr1) != 2 {
		return fmt.Errorf("bad email %s", email)
	}
	emailParsedDomain := arr1[1]
	val, ok := (*domainStat)[emailParsedDomain]
	if ok {
		(*domainStat)[emailParsedDomain] = val + 1
		return nil
	}
	arr2 := strings.Split(emailParsedDomain, ".")
	if len(arr2) != 2 {
		return fmt.Errorf("bad domain name %s", emailParsedDomain)
	}
	if lvl1 == arr2[1] {
		(*domainStat)[emailParsedDomain] = 1
	}
	return nil
}

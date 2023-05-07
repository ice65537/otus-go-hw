package hw10programoptimization

import (
	"bufio"
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
	if err == io.EOF {
		return nil
	}
	return err
}

func checkDomain(lvl1, email string, domainStat *DomainStat) error {
	arr1 := strings.Split(email, "@")
	if len(arr1) != 2 {
		return fmt.Errorf("bad email %s", email)
	}
	emailParsed_Domain := arr1[1]
	val, ok := (*domainStat)[emailParsed_Domain]
	if ok {
		(*domainStat)[emailParsed_Domain] = val + 1
		return nil
	}
	arr2 := strings.Split(emailParsed_Domain, ".")
	if len(arr2) != 2 {
		return fmt.Errorf("bad domain name %s", emailParsed_Domain)
	}
	if lvl1 == arr2[1] {
		(*domainStat)[emailParsed_Domain] = 1
	}
	return nil
}

package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
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
	result = make(DomainStat, 200)
	precs := parseRecords(r)
	for pr := range precs {
		if pr.err != nil {
			err = pr.err
			return
		}
		err = checkDomain(domain, pr.value.Email, &result)
		if err != nil {
			return
		}
	}
	return
}

type parsedRecord struct {
	value User
	err   error
}

func parseRecords(r io.Reader) (resultCh chan parsedRecord) {
	resultCh = make(chan parsedRecord, 100)

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
			err = easyjson.Unmarshal(uline, &user)
			if err != nil {
				break
			}
			user.Email = strings.ToLower(user.Email)
			resultCh <- parsedRecord{user, nil}
		}
		if err != io.EOF {
			resultCh <- parsedRecord{User{}, err}
		}
	}()
	return
}

func checkDomain(domain, email string, domainStat *DomainStat) error {
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
	if domain == arr2[1] {
		(*domainStat)[emailParsed_Domain] = 1
	}
	return nil
}

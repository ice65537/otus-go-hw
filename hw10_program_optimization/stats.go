package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	qjson "github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
)

type User struct {
	// ID       int
	// Name     string
	// Username string
	Email string
	// Phone    string
	// Password string
	// Address  string
}

type DomainStat map[string]int

var UseJSONBugerVsIter = true

func GetDomainStat(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	err = parseRecords(r, &result, domain)
	return
}

func parseRecords(r io.Reader, domainStat *DomainStat, lvl1 string) error {
	var uline []byte
	var user User
	var err error

	br := bufio.NewReaderSize(r, 1024)
	for err == nil {
		uline, err = br.ReadSlice('}')
		if err != nil {
			break
		}
		if UseJSONBugerVsIter {
			user.Email, err = qjson.GetString(uline, "Email")
			if err != nil {
				break
			}
		} else {
			user.Email = ""
			user.Email = jsoniter.Get(uline, "Email").ToString()
			if user.Email == "" {
				err = fmt.Errorf("email not found in %s", uline)
				break
			}
		}
		err = checkDomain(lvl1, user.Email, domainStat)
		if err != nil {
			break
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return err
}

func parseLowerEmail(email string) (lvl1, domain, lower string) {
	var dog, dot int
	lowerB := make([]byte, len(email))
	for i, v := range email {
		if v == '@' {
			dog = i
		}
		if v == '.' {
			dot = i
		}
		if v >= 'A' && v <= 'Z' {
			v += 'a' - 'A'
		}
		lowerB[i] = byte(v)
	}
	if dog == 0 || dot == 0 || dog > dot {
		return
	}
	lower = string(lowerB[:len(email)])
	lvl1 = lower[dot+1:]
	domain = lower[dog+1:]
	return
}

func checkDomain(lvl1, email string, domainStat *DomainStat) error {
	plvl1, pdomain, _ := parseLowerEmail(email)
	if plvl1 == "" {
		return fmt.Errorf("bad email %s", email)
	}
	val, ok := (*domainStat)[pdomain]
	if ok {
		(*domainStat)[pdomain] = val + 1
		return nil
	}
	if lvl1 == plvl1 {
		(*domainStat)[pdomain] = 1
	}
	return nil
}

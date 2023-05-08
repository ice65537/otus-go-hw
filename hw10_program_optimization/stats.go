package hw10programoptimization

import (
	"bufio"
	"bytes"
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

var (
	FullMemoryRead     bool = false
	UseJsonBugerVsIter bool = true
)

func GetDomainStat(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	if FullMemoryRead {
		err = parseRecords2(r, &result, domain)
	} else {
		err = parseRecords(r, &result, domain)
	}

	return
}

var ExpectedFileSize = 20_000_000

func parseRecords2(r io.Reader, domainStat *DomainStat, lvl1 string) error {
	var uline []byte
	var user User
	var err error

	data := make([]byte, ExpectedFileSize+1)

	n, err := io.ReadFull(r, data)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return err
	}
	if n > ExpectedFileSize {
		if data, err = io.ReadAll(r); err != nil {
			return err
		}
	}
	for _, uline = range bytes.Split(data, []byte("\n")) {
		if UseJsonBugerVsIter {
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
	return err
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
		if UseJsonBugerVsIter {
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
	lower = string(lowerB)
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

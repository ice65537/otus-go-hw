package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

func stringCheckLen(x string, info string) error {
	exactLength, err := strconv.Atoi(info)
	if err != nil {
		return err
	}
	if len(x) != exactLength {
		return ErrStrLen
	}
	return nil
}

func stringCheckRxp(x string, info string) error {
	rex, err := regexp.Compile(info)
	if err != nil {
		return err
	}
	if !rex.MatchString(x) {
		return ErrStrRxp
	}
	return nil
}

func stringCheckDict(x string, info string) error {
	var strSet map[string]struct{}
	strSetSlice := strings.Split(info, ",")
	strSet = make(map[string]struct{})
	for _, idx := range strSetSlice {
		strSet[idx] = struct{}{}
	}
	if _, ok := strSet[x]; !ok {
		return ErrStrNotFound
	}
	return nil
}

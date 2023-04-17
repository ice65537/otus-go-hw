package hw09structvalidator

import (
	"strconv"
	"strings"
)

func intCheckMin(x int, info string) error {
	min, err := strconv.Atoi(info)
	if err != nil {
		return err
	}
	if x < min {
		return ErrIntMin
	}
	return nil
}

func intCheckMax(x int, info string) error {
	max, err := strconv.Atoi(info)
	if err != nil {
		return err
	}
	if x > max {
		return ErrIntMax
	}
	return nil
}

func intCheckDict(x int, info string) error {
	var intSet map[int]struct{}
	intSetStr := strings.Split(info, ",")
	intSet = make(map[int]struct{})
	for _, v := range intSetStr {
		idx, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		intSet[idx] = struct{}{}
	}
	if _, ok := intSet[x]; !ok {
		return ErrIntNotFound
	}
	return nil
}

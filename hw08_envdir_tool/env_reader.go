package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envfiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envMap := make(Environment)

	for _, envfile := range envfiles {
		if strings.Contains(envfile.Name(), "=") {
			continue
		}
		f, err := os.Open(dir + "/" + envfile.Name())
		if err != nil {
			return nil, err
		}
		defer f.Close()

		fScanner := bufio.NewScanner(f)
		fScanner.Split(bufio.ScanLines)
		if fScanner.Scan() {
			sVal := strings.TrimRight(fScanner.Text(), " \t")
			sVal = strings.ReplaceAll(sVal, "\x00", "\n")
			envMap[envfile.Name()] = EnvValue{sVal, false}
		} else {
			envMap[envfile.Name()] = EnvValue{"", true}
		}
		f.Close()
	}
	return envMap, nil
}

package main

import (
	"fmt"
	"os"
)

func main() {
	envMap, err := ReadDir(os.Args[1])
	if err != nil {
		panic(fmt.Sprint("Ошибка чтения envdir: ", err))
	}

	retCode := RunCmd(os.Args[2:], envMap)
	os.Exit(retCode)
}

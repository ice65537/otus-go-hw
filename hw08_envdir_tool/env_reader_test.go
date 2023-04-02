package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	homedir := os.ExpandEnv("$HOME")
	check := func(e error) {
		if e != nil {
			fmt.Println(e)
		}
	}
	mkdir := func(dir string) {
		err := os.Mkdir(homedir+"/"+dir, 0o750)
		check(err)
	}
	clear := func(dir string) {
		os.RemoveAll(homedir + "/" + dir)
	}
	setenv := func(dir string, fname string, body string) {
		file, err := os.OpenFile(homedir+"/"+dir+"/"+fname, os.O_RDWR|os.O_CREATE, 0o664)
		check(err)
		_, err = file.WriteString(body)
		check(err)
		file.Close()
	}
	t.Run("TEST1. Простое чтение переменной из файла", func(t *testing.T) {
		mkdir("test1")
		defer clear("test1")
		setenv("test1", "ALPHA", "alpha_value")
		envMap, err := ReadDir(homedir + "/test1")
		check(err)
		require.Equal(t, "alpha_value", envMap["ALPHA"].Value, "Значение не прочиталось")
	})

	t.Run("TEST2. Чтение переменной из запрещенного файла", func(t *testing.T) {
		mkdir("test2")
		defer clear("test2")
		setenv("test2", "ALPHA=BETA", "alpha_value")
		envMap, err := ReadDir(homedir + "/test2")
		check(err)
		require.NotEqual(t, "alpha_value", envMap["ALPHA=BETA"].Value, "Значение прочиталось, а не должно было")
	})

	t.Run("TEST3. Чтение переменной из пустого файла", func(t *testing.T) {
		mkdir("test3")
		defer clear("test3")
		setenv("test3", "ALPHA", "")
		envMap, err := ReadDir(homedir + "/test3")
		check(err)
		require.True(t, envMap["ALPHA"].NeedRemove, "Параметр должен быть запланирован к удалению")
	})
}

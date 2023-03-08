package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported source file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var bufferSize int

	// Открытие файла
	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer source.Close()

	fi, err := source.Stat()
	if err != nil {
		return err
	}
	if fi.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if fi.Size() <= 0 {
		return ErrUnsupportedFile
	}
	_, err = source.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Открытие файла-копии
	target, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0o664)
	if err != nil {
		return err
	}
	defer target.Close()
	if limit == 0 {
		limit = fi.Size()
	}
	err = target.Truncate(0)
	if err != nil {
		return err
	}

	// Определение размера буфера для чтения
	bufferSize = 100
	osOut, err := exec.Command("stat", "--printf=%o", fromPath).Output()
	// osOut, err := exec.Command("ls", "-l").Output()
	if err != nil {
		fmt.Println(err)
	}
	i, err := strconv.Atoi(string(osOut))
	if err != nil {
		fmt.Println(err)
	} else {
		bufferSize = i
	}
	fmt.Printf("Memory buffer size: %d\r\n", bufferSize)

	// Копирование
	buffer := make([]byte, bufferSize)
	bar := pb.Full.Start64(limit)
	barWriter := bar.NewProxyWriter(target)
	for {
		numReaded, err := source.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		numToWrite := numReaded
		if int64(numReaded) > limit {
			numToWrite = int(limit)
		}
		if numToWrite > 0 {
			_, errw := barWriter.Write(buffer[:numToWrite])
			if errw != nil {
				return err
			}
		}
		limit -= int64(numToWrite)
		if limit <= 0 || err == io.EOF {
			break
		}
	}
	bar.Finish()
	return nil
}

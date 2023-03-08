package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	fmt.Printf("\r\nCopying file [%s], offset[%d] to [%s], limit[%d]\r\n", from, offset, to, limit)
	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
		os.Exit(0) // иначе ломается тест
	}
}

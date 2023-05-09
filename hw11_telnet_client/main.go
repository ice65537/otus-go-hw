package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	flagTimeout := flag.String("timeout", "10s", "timeout for server connect")
	flagHelp := flag.Bool("help", false, "view help message")
	flag.Parse()
	if *flagHelp {
		s := `
		go-telnet usage
		$ go-telnet --timeout=10s host port
		$ go-telnet mysite.ru 8080
		$ go-telnet --timeout=3s 1.1.1.1 123`
		fmt.Println(s)
		return
	}
	fmt.Println(flag.Args())
	timeout, err := time.ParseDuration(*flagTimeout)
	if err != nil {
		panic(fmt.Errorf("invalid timeout value [%s]", *flagTimeout))
	}
	host := flag.Args()[0]
	port := flag.Args()[1]
	client := NewTelnetClient(host+":"+port, timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	defer client.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Telnet client successfully connected to %s:%s\r\n", host, port)

	ctx, cancel := context.WithCancel(context.Background())
	var mtx sync.Mutex
	go receiver(ctx, &mtx)
	go sender(ctx, &mtx)
}

func receiver(ctx context.Context, mtx *sync.Mutex) {
	return
}

func sender(ctx context.Context, mtx *sync.Mutex) {
	return
}

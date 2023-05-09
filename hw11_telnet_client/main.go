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
	fmt.Fprintf(os.Stdout, "Telnet client successfully connected to %s:%s\r\n", host, port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mtx sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	go receiver(ctx, &mtx, &wg, client)
	go sender(ctx, &mtx, &wg, client)
	wg.Wait()
}

func receiver(ctx context.Context, mtx *sync.Mutex, wg *sync.WaitGroup, cli TelnetClient) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "Receive lock prepare\r\n")
			mtx.Lock()
			fmt.Fprintf(os.Stderr, "Receive lock get\r\n")
			err := cli.Receive()
			mtx.Unlock()
			fmt.Fprintf(os.Stderr, "Receive lock release\r\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Receive error: %s", err)
				cancel()
			}
		case <-ctx.Done():
			return
		}
	}
}

func sender(ctx context.Context, mtx *sync.Mutex, wg *sync.WaitGroup, cli TelnetClient) {
	defer wg.Done()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "Send lock prepare\r\n")
			mtx.Lock()
			fmt.Fprintf(os.Stderr, "Send lock get\r\n")
			err := cli.Send()
			mtx.Unlock()
			fmt.Fprintf(os.Stderr, "Send lock release\r\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Send error: %s", err)
				cancel()
			}
		case <-ctx.Done():
			return
		}
	}
}

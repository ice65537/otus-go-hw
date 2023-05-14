package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type ctxKey string

const keyCancel ctxKey = "Cancel"

func main() {
	flagTimeout := flag.String("timeout", "10s", "timeout for server connect")
	flagHelp := flag.Bool("help", false, "view help message")
	flagHello := flag.String("hello", "", "hello message for server")
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
		fmt.Fprintf(os.Stderr, "Invalid timeout value [%s]\n", *flagTimeout)
		return
	}
	host := flag.Args()[0]
	port := flag.Args()[1]
	scanBuffer := &bytes.Buffer{}
	scanBuffer.Grow(4096)
	client := NewTelnetClient(host+":"+port, timeout, io.NopCloser(scanBuffer), os.Stdout)
	err = client.Connect()
	defer client.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection error [%s]\n", err)
		return
	}
	fmt.Fprintf(os.Stderr, "Telnet client successfully connected to %s:%s\n", host, port)
	if *flagHello != "" {
		fmt.Fprint(scanBuffer, *flagHello+"\n")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, keyCancel, cancel)

	var wg sync.WaitGroup

	go scanner(ctx, scanBuffer)
	wg.Add(2)
	go receiver(ctx, client, &wg)
	go sender(ctx, client, &wg)
	wg.Wait()
	os.Exit(0)
}

func scanner(ctx context.Context, w io.Writer) {
	var inputBytes []byte
	var scanner *bufio.Scanner

	defer ctx.Value(keyCancel).(context.CancelFunc)()

	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Stdin get stat error: %s\n", err)
		return
	}
	if stdinStat.Mode()&os.ModeNamedPipe != 0 {
		fmt.Fprint(os.Stderr, "os.Stdin is in pipe mode\n")
		n, err := io.Copy(w, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Stdin pipe copy error: %s\n", err)
		}
		fmt.Fprintf(os.Stderr, "Stdin pipe copied %d bytes\n", n)
		return
	}

	scanner = bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		err = scanner.Err()
		inputBytes = append(scanner.Bytes(), '\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %s\n", err)
			return
		}
		if len(inputBytes) > 0 {
			_, err = w.Write(inputBytes)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Scan-write error: %s\n", err)
				return
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func receiver(ctx context.Context, cli TelnetClient, wg *sync.WaitGroup) {
	fmt.Fprint(os.Stderr, "Receiver started\n")
	defer wg.Done()
	defer ctx.Value(keyCancel).(context.CancelFunc)()
	for {
		err := cli.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Receive error: %s\n", err)
			return
		}
		select {
		default:
			continue
		case <-ctx.Done():
			fmt.Fprint(os.Stderr, "Receiver stopped\n")
			return
		}
	}
}

func sender(ctx context.Context, cli TelnetClient, wg *sync.WaitGroup) {
	fmt.Fprint(os.Stderr, "Sender started\n")
	defer wg.Done()
	defer ctx.Value(keyCancel).(context.CancelFunc)()
	for {
		err := cli.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Send error: %s\n", err)
			return
		}
		select {
		default:
			continue
		case <-ctx.Done():
			fmt.Fprint(os.Stderr, "Sender stopped\n")
			return
		}
	}
}

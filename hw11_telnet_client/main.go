package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type ctxKey string

const keyCancel ctxKey = "Cancel"

func main() {
	flagTimeout := flag.String("timeout", "10s", "timeout for server connect")
	flagHelp := flag.Bool("help", false, "view help message")
	flagHello := flag.String("hello", "WHATSUP", "hello message for server")
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
	scanBuffer := &bytes.Buffer{}
	in := io.NopCloser(scanBuffer)
	client := NewTelnetClient(host+":"+port, timeout, in, os.Stdout)
	err = client.Connect()
	defer client.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "Telnet client successfully connected to %s:%s\r\n", host, port)
	fmt.Fprintf(scanBuffer, *flagHello)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, keyCancel, cancel)

	go receiver(ctx, client)
	go sender(ctx, client)
	go scanner(ctx, scanBuffer)
	<-ctx.Done()
}

func scanner(ctx context.Context, w io.Writer) {
	defer ctx.Value(keyCancel).(context.CancelFunc)()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("[Send]>")
		if !scanner.Scan() {
			break
		}
		if _, err := w.Write(scanner.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "Scan-write error: %s", err)
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func receiver(ctx context.Context, cli TelnetClient) {
	defer ctx.Value(keyCancel).(context.CancelFunc)()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := cli.Receive()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Receive error: %s", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func sender(ctx context.Context, cli TelnetClient) {
	defer ctx.Value(keyCancel).(context.CancelFunc)()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := cli.Send()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Send error: %s", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

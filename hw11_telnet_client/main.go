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
	host := ""
	port := ""
	if len(flag.Args()) == 2 {
		host = flag.Args()[0]
		port = flag.Args()[1]
	} else {
		fmt.Fprintf(os.Stderr, "Invalid arg count [%d], expected 2 \n", len(flag.Args()))
		return
	}
	scanBuffer := &bytes.Buffer{}
	scanBuffer.Grow(4096)
	client := NewTelnetClient(host+":"+port, timeout, io.NopCloser(scanBuffer), os.Stdout)
	err = client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection error [%s]\n", err)
		return
	}
	defer client.Close()
	fmt.Fprintf(os.Stderr, "Telnet client successfully connected to %s:%s\n", host, port)
	if *flagHello != "" {
		fmt.Fprint(scanBuffer, *flagHello+"\n")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(3)
	go scanner(ctx, scanBuffer, &wg, cancel)
	go receiver(ctx, client, &wg, cancel)
	go sender(ctx, client, &wg, cancel)
	wg.Wait()
}

func scanner(ctx context.Context, w io.Writer, wg *sync.WaitGroup,
	parentCancel context.CancelFunc,
) {
	var inputBytes []byte
	var scanner *bufio.Scanner

	defer wg.Done()
	defer parentCancel()

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
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			err = scanner.Err()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Scan error: %s\n", err)
				return
			}
			inputBytes = append(scanner.Bytes(), '\n')
			if len(inputBytes) > 0 {
				if _, err = w.Write(inputBytes); err != nil {
					fmt.Fprintf(os.Stderr, "Scan-write error: %s\n", err)
					return
				}
			}
		}
	}
}

func receiver(ctx context.Context, cli TelnetClient, wg *sync.WaitGroup,
	parentCancel context.CancelFunc,
) {
	fmt.Fprint(os.Stderr, "Receiver started\n")
	defer wg.Done()
	defer parentCancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Fprint(os.Stderr, "Receiver stopped\n")
			return
		default:
			if err := cli.Receive(); err != nil {
				fmt.Fprintf(os.Stderr, "Receive error: %s\n", err)
				return
			}
		}
	}
}

func sender(ctx context.Context, cli TelnetClient, wg *sync.WaitGroup,
	parentCancel context.CancelFunc,
) {
	fmt.Fprint(os.Stderr, "Sender started\n")
	defer wg.Done()
	defer parentCancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Fprint(os.Stderr, "Sender stopped\n")
			return
		default:
			err := cli.Send()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Send error: %s\n", err)
				return
			}
		}
	}
}
